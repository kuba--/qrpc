package qrpc

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/kuba--/qrpc/internal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// cluster structure aggregates info about self (unique key + local address),
// and lookup table for other peers.
type cluster struct {
	key     string
	laddr   string
	timeout time.Duration

	sync.RWMutex
	peers map[string]*internal.Peer
}

// join implements kind of cluster discovery mechanism.
// It uses gossip protocol for internal communication.
// When it's done the Cluster.peers lookup table should be set.
func (c *cluster) Join(addrs []string) {
	var peers []*internal.Peer

	for _, addr := range addrs {
		if c.laddr == addr {
			continue
		}
		peers = append(peers, &internal.Peer{Addr: addr})
	}

	if len(peers) > 0 {
		c.joinPeers(peers)
	}
}

// joinPeers sends internal.JoinRequest(s) to all discovered peers.
func (c *cluster) joinPeers(peers []*internal.Peer) {
	for _, peer := range peers {
		if _, exists := c.getPeer(peer.Key); exists || c.laddr == peer.Addr {
			continue
		}

		resp, err := c.request(context.Background(), peer.Addr, &internal.JoinRequest{Key: c.key, Laddr: c.laddr})
		if err != nil {
			log.Println(err)
			continue
		}
		joinresp := resp.(*internal.JoinResponse)

		c.setPeer(joinresp.Key, peer.Addr)
		c.joinPeers(joinresp.Peers)
	}
}

// unjoin sends (fire & forget) internal.UnjoinRequest(s) to all healthy peers.
// When it's done, we should have standalone running server.
func (c *cluster) unjoin() {
	for _, peer := range c.copyPeers() {
		if c.key == peer.Key || c.laddr == peer.Addr {
			continue
		}

		go func(addr string) {
			if _, err := c.request(context.Background(), addr, &internal.UnjoinRequest{Key: c.key}); err != nil {
				log.Println(err)
			}
		}(peer.Addr)
	}
}

// watch runs watchdog timer forever (or if we received ctx.Done())
// Every interval duration internal.PingRequest is sent to all peers.
// We can cancel "watching" by calling cancel function associated with context.
func (c *cluster) watch(ctx context.Context, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Watch", ctx.Err())
			return

		case <-time.After(interval):
			wg := &sync.WaitGroup{}
			for _, peer := range c.copyPeers() {
				wg.Add(1)
				go func(key, addr string) {
					defer wg.Done()
					if _, err := c.request(ctx, addr, &internal.PingRequest{Key: c.key, Laddr: c.laddr}); err != nil {
						c.deletePeer(key)
						log.Println(addr, err)
					} else {
						c.setPeer(key, addr)
					}
				}(peer.Key, peer.Addr)
			}
			wg.Wait()
		}
	}
}

// getPeer, if peer for a given key exists, the method returns pointer to the internal.Peer, and true value.
// Otherwise (if peer does not exists in cluster's lookup table), nil and false.
func (c *cluster) getPeer(key string) (*internal.Peer, bool) {
	c.RLock()
	defer c.RUnlock()

	peer, exists := c.peers[key]
	return peer, exists
}

// setPeer sets/adds a peer in cluster's lookup table.
// If key already exists, address will be overwritten.
func (c *cluster) setPeer(key, addr string) {
	c.Lock()
	defer c.Unlock()
	c.peers[key] = &internal.Peer{Key: key, Addr: addr}
}

// deletePeer deletes a peer for a given key from cluster's lookup table.
func (c *cluster) deletePeer(key string) {
	c.Lock()
	delete(c.peers, key)
	c.Unlock()
}

// copyPeers returns a copy of all known and healthy peers (except specified in excluded list).
func (c *cluster) copyPeers(excludeKeys ...string) []*internal.Peer {
	var m []*internal.Peer

	c.RLock()
NextPeer:
	for _, peer := range c.peers {
		for _, key := range excludeKeys {
			if key == peer.Key {
				continue NextPeer
			}
		}
		m = append(m, &internal.Peer{Key: peer.Key, Addr: peer.Addr})
	}
	c.RUnlock()

	return m
}

// write replicates a message across all known cluster's peers.
func (c *cluster) write(key string, msg []byte) {
	log.Printf("[%s] -> write: %s\n", c.laddr, key)
	wg := &sync.WaitGroup{}
	for _, peer := range c.copyPeers() {
		wg.Add(1)
		go func(addr string) {
			if _, err := c.request(context.Background(), addr, &internal.WriteRequest{Key: key, Msg: msg}); err != nil {
				log.Println(err)
			}
			wg.Done()
		}(peer.Addr)
	}
	wg.Wait()
}

// erase deletes already read message from all known cluster's peers.
func (c *cluster) erase(key string) {
	log.Printf("[%s] -> erase: %s\n", c.laddr, key)
	wg := &sync.WaitGroup{}
	for _, peer := range c.copyPeers() {
		wg.Add(1)
		go func(addr string) {
			if _, err := c.request(context.Background(), addr, &internal.EraseRequest{Key: key}); err != nil {
				log.Println(c.laddr, err)
			}
			wg.Done()
		}(peer.Addr)
	}
	wg.Wait()
}

// request is a generic method to send internal requests to qrpc servers.
// request internally uses a copy of parent context with timeout (see: Cluster.timeout)
func (c *cluster) request(parent context.Context, addr string, req interface{}) (interface{}, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, _ := context.WithTimeout(parent, c.timeout)
	switch req.(type) {
	case *internal.JoinRequest:
		return internal.NewGossipClient(conn).Join(ctx, req.(*internal.JoinRequest))

	case *internal.UnjoinRequest:
		return internal.NewGossipClient(conn).Unjoin(ctx, req.(*internal.UnjoinRequest))

	case *internal.PingRequest:
		return internal.NewGossipClient(conn).Ping(ctx, req.(*internal.PingRequest))

	case *internal.WriteRequest:
		return internal.NewGossipClient(conn).Write(ctx, req.(*internal.WriteRequest))

	case *internal.EraseRequest:
		return internal.NewGossipClient(conn).Erase(ctx, req.(*internal.EraseRequest))
	}

	return nil, errors.New("invalid request")
}
