package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
	"testing"
)

// this test is pretty bs, needs improvements
func TestRenderFrame(t *testing.T) {
	ctx := context.Background()
	const bufSize = 1024 * 1024
	var listener *bufconn.Listener
	listener = bufconn.Listen(bufSize)
	// minio client data
	endpoint := "localhost:9000"
	accessKeyID := "minioaccesskeyid"
	secretAccessKey := "miniosecretaccesskey"
	useSSL := false

	// Initialize minio client object.
	mC, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		t.Log("minioClient oopsie:", err)
	}
	//log.Println("mC:", minioClient)
	minioClient = mC

	gcsCacheDir = os.TempDir()

	srv := grpc.NewServer()
	pb.RegisterRenderServer(srv, server{})
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
	client := pb.NewRenderClient(conn)
	//renderRequest := /* TODO: fill this in */
	renderRequest := &pb.RenderRequest{
		GcsOutputBase: "gifbucket",
		ObjPath:       "airboat.obj",
		Assets:        []string{},
		Rotation:      2.0,
		Iterations:    1,
	}
	res, err := client.RenderFrame(ctx, renderRequest)
	if err != nil {
		t.Log("error rendering frame:", err)
	}
	log.Println("response:", res)
}
