package api

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// request is a generic function to send api requests to qrpc servers.
func Request(ctx context.Context, addr string, req interface{}) (interface{}, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	switch req.(type) {
	case *SendRequest:
		return NewQRPCClient(conn).Send(ctx, req.(*SendRequest))

	case *ReceiveRequest:
		return NewQRPCClient(conn).Receive(ctx, req.(*ReceiveRequest))
	}

	return nil, errors.New("invalid request")
}
