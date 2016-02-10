package qrpc

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kuba--/qrpc/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestClusterJoin(t *testing.T) {
	s := newCluster(3)
	errchan0 := s[0].Start("9000")
	errchan1 := s[1].Start("9001", ":9000")
	errchan2 := s[2].Start("9002", ":9001")
	defer stopCluster(s)

	select {
	case err0 := <-errchan0:
		t.Fatal("err0", err0)
	case err1 := <-errchan1:
		t.Fatal("err1", err1)
	case err2 := <-errchan2:
		t.Fatal("err2", err2)

	case <-time.After(5 * time.Second):
		if _, exists := s[2].cluster.getPeer(s[0].cluster.key); !exists {
			t.Error(s[2].cluster.copyPeers())
		}
		if _, exists := s[2].cluster.getPeer(s[1].cluster.key); !exists {
			t.Error(s[2].cluster.copyPeers())
		}

		if _, exists := s[1].cluster.getPeer(s[0].cluster.key); !exists {
			t.Error(s[1].cluster.copyPeers())
		}
		if _, exists := s[1].cluster.getPeer(s[2].cluster.key); !exists {
			t.Error(s[1].cluster.copyPeers())
		}

		if _, exists := s[0].cluster.getPeer(s[1].cluster.key); !exists {
			t.Error(s[0].cluster.copyPeers())
		}
		if _, exists := s[0].cluster.getPeer(s[2].cluster.key); !exists {
			t.Error(s[0].cluster.copyPeers())
		}
	}
}

func TestClusterStop(t *testing.T) {
	s := newCluster(3)
	errchan0 := s[0].Start("8000")
	errchan1 := s[1].Start("8001", ":8000")
	errchan2 := s[2].Start("8002", ":8001")
	defer stopCluster([]*Server{s[2]})

	select {
	case err0 := <-errchan0:
		t.Fatal("err0", err0)
	case err1 := <-errchan1:
		t.Fatal("err1", err1)
	case err2 := <-errchan2:
		t.Fatal("err2", err2)

	case <-time.After(time.Second):
		go s[0].Stop()
		go s[1].Stop()
	}

	select {
	case err2 := <-errchan2:
		t.Fatal("err2", err2)

	case <-time.After(time.Second):
		if _, exists := s[2].cluster.getPeer(s[0].cluster.key); exists {
			t.Error(s[2].cluster.copyPeers())
		}
		if _, exists := s[2].cluster.getPeer(s[1].cluster.key); exists {
			t.Error(s[2].cluster.copyPeers())
		}
	}
}

func TestClusterSendReceive(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lookup := make(map[string]string)

	s := newCluster(3)
	errchan0 := s[0].Start("7000")
	errchan1 := s[1].Start("7001", ":7000")
	errchan2 := s[2].Start("7002", ":7000")
	defer stopCluster(s)

	select {
	case err0 := <-errchan0:
		t.Fatal("err0", err0)
	case err1 := <-errchan1:
		t.Fatal("err1", err1)
	case err2 := <-errchan2:
		t.Fatal("err2", err2)
	case <-time.After(time.Second):
		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)

			str := strconv.Itoa(i)
			go func(topic string) {
				defer wg.Done()

				_, err := s[rand.Intn(3)].Send(context.TODO(), &api.SendRequest{Topic: topic, Msg: []byte(topic)})
				if err != nil {
					t.Error(err)
				}

				resp, err := s[rand.Intn(3)].Receive(context.TODO(), &api.ReceiveRequest{Topic: topic})
				if err != nil {
					if grpc.Code(err) == codes.NotFound {
						log.Println(err)
					} else {
						t.Error(err)
					}
				} else {
					if v, b := lookup[resp.Key]; b {
						t.Error("Msg: %v already exists: %v\n", resp, v)
					}
					lookup[resp.Key] = string(resp.Msg)
					log.Println(resp)
				}
			}(str)
		}
		wg.Wait()
	}
}

func newCluster(n int) (servers []*Server) {
	for i := 0; i < n; i++ {
		servers = append(servers, NewServer("/tmp/qrpc-"+strconv.Itoa(i), 1024))
	}
	return
}

func stopCluster(servers []*Server) {
	for _, s := range servers {
		s.Stop()
	}
}
