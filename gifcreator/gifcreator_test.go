package main

import (
	"context"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"gopkg.in/redis.v5" //very outdated api version
	"log"
)

func TestStartJob(t *testing.T) {
	ctx := context.Background()
	const bufSize = 1024 * 1024
	var listener *bufconn.Listener
	listener = bufconn.Listen(bufSize)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

	client := pb.NewGifCreatorClient(conn)
	startJobRequest := &pb.StartJobRequest{
		Name: "katiething",
		ProductToPlug: 1,
	}
	res, err := client.StartJob(ctx, startJobRequest)
	if err != nil {
		t.Log("error starting job:", err)
	}
	log.Println("response:", res)

}
