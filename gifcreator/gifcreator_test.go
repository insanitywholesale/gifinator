package main

import (
	"context"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

func TestStartJob(t *testing.T) {
	ctx := context.Background()
	const bufSize = 1024 * 1024
	var listener *bufconn.Listener
	listener = bufconn.Listen(bufSize)

	srv := grpc.NewServer()
	pb.RegisterGifCreatorServer(srv, server{})
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			t.Log("server error:", err)
		}
	}()

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Log("grpc dial error:", err)
	}
	defer conn.Close()

	// add database connection stuff
}
