package qrpc

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/kuba--/qrpc/api"
)

const (
	topic = "qrpc-test"
	data  = "/tmp/qrpc-test"
	cache = 1024
	port  = "9090"
)

func TestServer(t *testing.T) {
	s := NewServer(data, cache)
	errchan := s.Start(port)
	defer s.Stop()

	select {
	case err := <-errchan:
		t.Error(err)
	default:
		testServerSend(t)
		testServerReceive(t)
	}
}

func testServerSend(t *testing.T) {
	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(cnt int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			msg := strconv.Itoa(cnt)
			_, err := api.Request(ctx, ":"+port, &api.SendRequest{topic, []byte(msg)})
			if err != nil {
				cancel()
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
}

func testServerReceive(t *testing.T) {
	wg := &sync.WaitGroup{}
	lookup := make(map[string]string)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			if r, err := api.Request(ctx, ":"+port, &api.ReceiveRequest{topic}); err != nil {
				if grpc.Code(err) == codes.NotFound {
					return
				}
				cancel()
				t.Error(err)
			} else {
				fmt.Println(r)
				resp := r.(*api.ReceiveResponse)
				if v, b := lookup[resp.Key]; b {
					t.Errorf("Msg: %v already exists: %v\n", resp, v)
				}
				lookup[resp.Key] = string(resp.Msg)
			}
		}()
	}
	wg.Wait()
}
