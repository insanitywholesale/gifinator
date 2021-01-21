/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/fogleman/pt/pt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	//"strings"
)

type server struct{}

var (
	gcsCacheDir string
	minioClient *minio.Client
)

func cacheMinioObjToDisk(ctx context.Context, fileName string) (string, error) {
	basePath := "/tmp/objcache"
	fullPath := basePath + "/" + fileName
	log.Println("filename", fileName)
	err := minioClient.FGetObject(context.Background(), "gifbucket", fileName, fullPath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("err:", err)
		return "", err
	}
	return fullPath, nil
}

func renderImage(objectPath string, rotation float64, iterations int32) (string, error) {
	scene := pt.Scene{}

	// create materials
	objMat := pt.GlossyMaterial(pt.Black, 1.2, pt.Radians(30))
	wall := pt.GlossyMaterial(pt.HexColor(0xFCFAE1), 1.5, pt.Radians(10))
	light := pt.LightMaterial(pt.White, 80)

	// add walls and lights
	scene.Add(pt.NewCube(pt.V(-10, -1, -10), pt.V(-2, 10, 10), wall))
	scene.Add(pt.NewCube(pt.V(-10, -1, -10), pt.V(10, 0, 10), wall))
	scene.Add(pt.NewSphere(pt.V(4, 10, 1), 1, light))

	// load and transform gopher mesh
	mesh, err := pt.LoadOBJ(objectPath, objMat)
	if err != nil {
		return "", err
	}

	mesh.Transform(pt.Rotate(pt.V(0, 1, 0), pt.Radians(-10)))
	mesh.SmoothNormals()
	mesh.FitInside(pt.Box{pt.V(-1, 0, -1), pt.V(1, 2, 1)}, pt.V(0.5, 0, 0.5))
	mesh.Transform(pt.Rotate(pt.V(0, 1, 0), pt.Radians(rotation)))
	scene.Add(mesh)

	// position camera
	camera := pt.LookAt(pt.V(4, 1, 0), pt.V(0, 0.9, 0), pt.V(0, 1, 0), 30)

	// render the scene
	sampler := pt.NewSampler(16, 16)
	renderer := pt.NewRenderer(&scene, &camera, sampler, 300, 300)

	// TODO(jessup) Fix this for better entropy
	imagePath := os.TempDir() + "/final_img_itr_%d_" + strconv.FormatInt(int64(rand.Intn(10000)), 16) + ".png"
	renderer.IterativeRender(imagePath, int(iterations))

	return fmt.Sprintf(imagePath, iterations), nil
}

func (server) RenderFrame(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	fmt.Fprintf(os.Stdout, "starting render job - object: %s, angle: %f\n", req.ObjPath, req.Rotation)

	// Load main object (.obj) file
	// TODO: fix nil pointer dereference happenning here
	//objGcsObj, err := minioClient.GetObject(ctx, "gifbucket", req.ObjPath, minio.GetObjectOptions{})
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "error getting object %s, err: %v\n", req.ObjPath, err)
	//}
	objFilepath, err := cacheMinioObjToDisk(ctx, req.ObjPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error caching %s, err: %v\n", req.ObjPath, err)
		return nil, err
	}

	// Load the assets
	for _, element := range req.Assets {
		assetGcsObj, err := minioClient.GetObject(ctx, "gifbucket", element /*maybe*/, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println("err caching object:", assetGcsObj, "error:", err)
			return nil, err
		}
	}

	// Create and render a scene seeded with the object we loaded
	fmt.Fprintf(os.Stdout, "starting actual render - object: %s, angle: %f\n", req.ObjPath, req.Rotation)
	imgPath, err := renderImage(objFilepath, float64(req.Rotation), req.Iterations)
	if err != nil {
		log.Println("renderImage error:", err)
	}
	fmt.Fprintf(os.Stdout, "finished actual render - object: %s, angle: %f\n", req.ObjPath, req.Rotation)

	// name of the final image that will be saved to minio
	imgFinalName := fmt.Sprintf("%s.image_%.0frad.png", req.GcsOutputBase, req.Rotation)
	// change renderImage to return a filename to use here so I don't have to do the following
	// find last occurance of / and get the string from there to the end
	// the slice trick results in an error
	// I winged it and didn't think about what happens when it returns -1
	// AKA what happens when the imgPath is not GCS-related

	//imgFileName := strings.TrimLeft(imgPath[strings.LastIndex(imgPath, "/"):], "/")
	uploadInfo, err := minioClient.FPutObject(ctx, "gifbucket", imgFinalName, imgPath, minio.PutObjectOptions{})
	if err != nil {
		log.Println("error uploading image to minio", err)
	}
	fmt.Println("uploaded:", uploadInfo)
	response := pb.RenderResponse{GcsOutput: imgFinalName}
	return &response, nil
}

func main() {
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
		log.Fatalln(err)
	}
	// SOMEHOW this line fixes the nil pointer dereference issue
	minioClient = mC

	// lines intentionally left blank

	serving_port := os.Getenv("RENDER_PORT")
	if serving_port == "" {
		serving_port = "8080"
	}
	i, err := strconv.Atoi(serving_port)
	if (err != nil) || (i < 1) {
		log.Fatalf("please set env var RENDER_PORT to a valid port")
		return
	}

	l, err := net.Listen("tcp", ":"+serving_port)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
		return
	}
	log.Println("running on port:", serving_port)

	gcsCacheDir = os.TempDir()

	srv := grpc.NewServer()
	pb.RegisterRenderServer(srv, server{})
	srv.Serve(l)
}
