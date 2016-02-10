package qrpc

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"

	"diskv"

	"github.com/kuba--/qrpc/api"
	"github.com/kuba--/qrpc/internal"
)

var (
	// Confuse `go vet' to not check this `Errorf' call.
	// See https://github.com/grpc/grpc-go/issues/90
	grpcErrorf = grpc.Errorf
)

// Server structure represents QRPC node.
type Server struct {
	rpc           *grpc.Server
	cluster       *cluster
	cancelWatcher context.CancelFunc

	mtx  *sync.Mutex
	data *diskv.Diskv
}

// NewServer creates a qrpc server which has not started to accept requests yet.
func NewServer(dataDir string, cacheSize uint64) *Server {
	return &Server{
		rpc: grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32)),
		cluster: &cluster{
			key:     fmt.Sprintf("peer-%s", keySuffix()),
			timeout: 3 * time.Second,
			peers:   make(map[string]*internal.Peer),
		},
		cancelWatcher: nil,

		mtx: &sync.Mutex{},
		data: diskv.New(diskv.Options{
			BasePath:     dataDir,
			Transform:    transformKey,
			CacheSizeMax: cacheSize,
		}),
	}
}

// Start starts a qrpc server on specified port and tries to connect to existing cluster (if peers were specified).
// Returns an error channel which will be populated when server fail.
func (s *Server) Start(port string, peers ...string) <-chan error {
	log.Printf("-> Start[%s](%s, %d): %s\t%v", s.cluster.key, s.data.BasePath, s.data.CacheSizeMax, port, peers)

	errchan := make(chan error, 1)

	api.RegisterQRPCServer(s.rpc, s)
	internal.RegisterGossipServer(s.rpc, s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		errchan <- err
		return errchan
	}
	go func() { errchan <- s.rpc.Serve(lis) }()

	s.cluster.laddr = lis.Addr().String()
	s.cluster.Join(peers)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelWatcher = cancel

	go s.cluster.watch(ctx, time.Second)

	return errchan
}

// Stop stops a qrpc server.
func (s *Server) Stop() {
	log.Println("-> Stop")
	if s.cancelWatcher != nil {
		s.cancelWatcher()
	}
	s.cluster.unjoin()

	s.rpc.Stop()
}

// Send implements 'send request endpoint'.
// Data will be stored in diskv and replicated across cluster's peers.
func (s *Server) Send(ctx context.Context, req *api.SendRequest) (*api.SendResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		key := newKey(req.Topic)
		if err := s.data.Write(key, req.Msg); err != nil {
			return nil, err
		}

		s.cluster.write(key, req.Msg)
		return &api.SendResponse{key}, nil
	}
}

// Receive implements 'receive request endpoint'.
// Data will be read and erased from diskv.
// The server will also try to erase the read key from cluster's peers.
func (s *Server) Receive(ctx context.Context, req *api.ReceiveRequest) (*api.ReceiveResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		var (
			key string
			msg []byte
			err error
		)
		cancel := make(chan struct{}, 1)

		s.mtx.Lock()
		for key = range s.data.KeysPrefix(keyPrefix(req.Topic), cancel) {
			msg, err = s.data.Read(key)
			if err != nil {
				log.Println(err)
				continue
			}

			if err = s.data.Erase(key); err != nil {
				log.Println(err)
				continue
			}

			// close KeysPrefix channel and break the loop
			cancel <- struct{}{}
			break
		}
		s.mtx.Unlock()

		if err == nil && key != "" {
			s.cluster.erase(key)
			return &api.ReceiveResponse{key, msg}, nil
		}
	}
	return nil, grpcErrorf(codes.NotFound, "data not found (%s)", req.Topic)
}

// Join implements 'internal join request endpoint'. Updates address for requesting peer.
// Returns list of already known peers.
func (s *Server) Join(ctx context.Context, req *internal.JoinRequest) (*internal.JoinResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		hostport, err := peerFromContext(ctx, req.Laddr)
		if err != nil {
			return nil, err
		}

		peers := s.cluster.copyPeers(req.Key)
		log.Println("Join (", s.cluster.key, s.cluster.laddr, ") <- (", req.Key, hostport, ") :", peers)

		s.cluster.setPeer(req.Key, hostport)
		return &internal.JoinResponse{Key: s.cluster.key, Peers: peers}, nil
	}
}

// Unjoin implements 'internal unjoin request endpoint'.
// Deletes requesting peer from lookup table.
func (s *Server) Unjoin(ctx context.Context, req *internal.UnjoinRequest) (*internal.UnjoinResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		s.cluster.deletePeer(req.Key)
		return &internal.UnjoinResponse{}, nil
	}
}

// Ping implements 'internal ping request endpoint'.
// The endpoint is mostly used by cluster's watcher for service discovery.
func (s *Server) Ping(ctx context.Context, req *internal.PingRequest) (*internal.PingResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		hostport, err := peerFromContext(ctx, req.Laddr)
		if err != nil {
			return nil, err
		}

		s.cluster.setPeer(req.Key, hostport)
		return &internal.PingResponse{}, nil
	}
}

// Write implements 'internal write request endpoint'.
// The endpoint is mostly used to replicate data across cluster's peers.
// If requesting req.Key already exists then data will be overwritten.
func (s *Server) Write(ctx context.Context, req *internal.WriteRequest) (*internal.WriteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		err := s.data.Write(req.Key, req.Msg)
		if err != nil {
			return nil, err
		}
		return &internal.WriteResponse{req.Key}, nil
	}
}

// Erase implements 'internal erase request endpoint'.
// The endpoint is used to delete the key from the server (becuase data were read from other peer).
func (s *Server) Erase(ctx context.Context, req *internal.EraseRequest) (*internal.EraseResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		s.mtx.Lock()
		defer s.mtx.Unlock()

		if err := s.data.Erase(req.Key); err != nil {
			log.Println(err)
			return nil, err
		}

		return &internal.EraseResponse{req.Key}, nil
	}
}

// peerFromContext extracts host from the context and port the from local address.
// Returns an address based on extracted host:port
func peerFromContext(ctx context.Context, laddr string) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", errors.New("peer information in context does not exist")
	}

	host, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return "", err
	}

	_, port, err := net.SplitHostPort(laddr)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(host, port), nil
}
