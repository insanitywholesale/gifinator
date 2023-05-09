package main

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/alicebob/miniredis" // unused for now
	"github.com/go-redis/redis/v8"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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
func TestStartJob(t *testing.T) {
	// set up base variables
	minioName := os.Getenv("MINIO_NAME")
	if minioName == "" {
		minioName = "localhost"
	}
	minioPort := os.Getenv("MINIO_PORT")
	if minioPort == "" {
		minioPort = "9000"
	}
	endpoint = minioName + ":" + minioPort
	minioBucket = "gifbucket"
	accessKeyID = "minioaccesskeyid"
	secretAccessKey = "miniosecretaccesskey"
	scenePath = "/tmp/scene"
	ctx := context.Background()
	const bufSize = 1024 * 1024
	listener := bufconn.Listen(bufSize)

	redisName := os.Getenv("REDIS_NAME")
	if redisName == "" {
		redisName = "localhost"
	}
	t.Log(redisName)
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	t.Log(redisPort)
	// initialize redis client
	// mr, _ := miniredis.Run()
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisName + ":" + redisPort,
		// Addr:     mr.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if redisClient == nil {
		t.Error("redis client is nil")
	}

	// make new grpc server
	srv := grpc.NewServer()
	pb.RegisterGifCreatorServer(srv, server{})
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			t.Error("server error:", err)
		}
	}()

	// dial above grpc server
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Error("grpc dial error:", err)
	}
	defer conn.Close()

	// create new client in order to run StartJob
	client := pb.NewGifCreatorClient(conn)
	// create a StartJobRequest to use
	startJobRequest := &pb.StartJobRequest{
		Name:          "k8s",
		ProductToPlug: 2,
	}
	// run StartJob with the above request
	res, err := client.StartJob(ctx, startJobRequest)
	if err != nil {
		t.Fatal("error starting job:", err)
	}
	t.Log("response:", res)
	t.Log("response id:", res.JobId)
}

func TestWorkerMode(t *testing.T) {
	// set up base variables
	*workerMode = true
	renderHostAddr := "localhost:8080"
	scenePath = "/scene"

	// initialize redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// dial grpc server started by the renderer
	conn, err := grpc.Dial(renderHostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Log("error connecting to renderer at", renderHostAddr, err)
	}

	// create new render client with the above connection
	renderClient = pb.NewRenderClient(conn)

	log.Println("connected to renderer at", renderHostAddr)

	// poll the task queue and lease tasks
	err = leaseNextTask()
	if err != nil {
		t.Error("error working on task", err)
	}
	time.Sleep(10 * time.Millisecond)
	conn.Close()
	// needs additions to actually test
}
