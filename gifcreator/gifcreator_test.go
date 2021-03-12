package main

import (
	"context"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	//"gopkg.in/redis.v5" //very outdated api version
	"github.com/go-redis/redis/v8"
	"log"
	"net"
	"testing"
	"time"
)

// TODO: tests b bork; two different instances (one in server mode, one in worker mode) are needed
/*
func TestCompileGifs(t *testing.T) {
	prefix := "out."
	tCtx := context.Background()
	link, err := compileGifs(prefix, tCtx)
	if err != nil {
		t.Fatal("failed to compile gif", err)
	}
	log.Println("link:", link)
}
*/

// scenePath, redisClient and server should be changed
// in order to be able to run go test gifcreator/gifcreator_test.go
// from root of repo
// TODO: this test errors out, need to fix renderer or leaseNextTask or both
func TestStartJob(t *testing.T) {

	// set up base variables
	endpoint = "localhost:9000"
	minioBucket = "gifbucket"
	accessKeyID = "minioaccesskeyid"
	secretAccessKey = "miniosecretaccesskey"
	scenePath = "/scene"
	ctx := context.Background()
	const bufSize = 1024 * 1024
	var listener *bufconn.Listener
	listener = bufconn.Listen(bufSize)

	// initialize redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// make new grpc server
	srv := grpc.NewServer()
	pb.RegisterGifCreatorServer(srv, server{})
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			t.Log("server error:", err)
		}
	}()

	// dial above grpc server
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

	//create new client in order to run StartJob
	client := pb.NewGifCreatorClient(conn)
	//create a StartJobRequest to use
	startJobRequest := &pb.StartJobRequest{
		Name:          "k8s",
		ProductToPlug: 2,
	}
	//run StartJob with the above request
	res, err := client.StartJob(redisContext, startJobRequest)
	if err != nil {
		t.Log("error starting job:", err)
	}
	log.Println("response:", res)
	log.Println("response id:", res.JobId)

}

func TestWorkerMode(t *testing.T) {
	// set up base variables
	*workerMode = true
	renderHostAddr := "localhost:8080"
	scenePath = "/tmp/scene"

	// initialize redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// dial grpc server started by the renderer
	conn, err := grpc.Dial(renderHostAddr, grpc.WithInsecure())
	if err != nil {
		t.Log("error connecting to renderer at", renderHostAddr, err)
	}

	// create new render client with the above connection
	renderClient = pb.NewRenderClient(conn)

	log.Println("connected to renderer at", renderHostAddr)

	// poll the task queue and lease tasks
	err = leaseNextTask()
	if err != nil {
		t.Log("error working on task", err)
	}
	time.Sleep(10 * time.Millisecond)
	conn.Close()
	// needs additions to actually test
}
