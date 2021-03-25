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
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang/freetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"golang.org/x/image/font/gofont/gobold"
	"google.golang.org/grpc"
	"image"
	"image/gif"
	"image/png"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type server struct{}

type renderJob struct {
	Status         pb.GetJobResponse_Status
	FinalImagePath string
}

type renderTask struct {
	Frame       int64
	Caption     string
	ProductType pb.Product
}

var (
	redisClient     *redis.Client
	renderClient    pb.RenderClient
	scenePath       string
	deploymentId    string
	workerMode      = flag.Bool("worker", false, "run in worker mode rather than server")
	redisName = "localhost"
	redisPort = "6379"
	minioBucket     string
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          = false
)

func transform(inputPath string, jobId string) (bytes.Buffer, error) {
	var transformed bytes.Buffer
	tmpl, err := template.ParseFiles(inputPath)
	if err != nil {
		return transformed, err
	}
	err = tmpl.Execute(&transformed, jobId)
	if err != nil {
		return transformed, err
	}
	return transformed, nil
}

// Utility function to upload something to minio
func upload(outBytes []byte, outputPath string, mimeType string, client *minio.Client, ctx context.Context) error {
	objName := outputPath
	uploadInfo, err := client.PutObject(
		ctx,
		minioBucket,
		objName,
		bytes.NewReader(outBytes),
		int64(len(outBytes)),
		minio.PutObjectOptions{ContentType: mimeType},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Add the text given to the badge image
func addLabel(img *image.NRGBA, x, y int, label string) error {
	fontSize := float64(120)
	f, err := freetype.ParseFont(gobold.TTF)
	if err != nil {
		return err
	}
	c := freetype.NewContext()
	c.SetDPI(float64(72))
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(img)
	pt := freetype.Pt(x, y+int(c.PointToFixed(fontSize)>>6))
	_, err = c.DrawString(label, pt)
	return err
}

func (server) StartJob(ctx context.Context, req *pb.StartJobRequest) (*pb.StartJobResponse, error) {
	redisContext := context.Background()
	// Retrieive the next job ID from Redis
	jobId, err := redisClient.Incr(redisContext, "gifjob_counter").Result()
	if err != nil {
		return nil, err
	}
	jobIdStr := strconv.FormatInt(jobId, 10)

	// Create a new RenderJob queue for that job
	var job = renderJob{
		Status: pb.GetJobResponse_PENDING,
	}
	payload, _ := json.Marshal(job)
	err = redisClient.Set(redisContext, "job_gifjob_"+jobIdStr, payload, 0).Err()
	if err != nil {
		return nil, err
	}

	// Make a new minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println("error making minio client", err)
		return nil, err
	}

	// Set what mascot will be used
	var productString string
	switch req.ProductToPlug {
	case pb.Product_GRPC:
		productString = "grpc"
		break
	case pb.Product_KUBERNETES:
		productString = "k8s"
		break
	default:
		productString = "gopher"
	}

	// Generate the assets needed to render the frame, and upload them to minio
	t, err := transform(scenePath+"/"+productString+".obj.tmpl", jobIdStr)
	if err != nil {
		return nil, err
	}
	err = upload(t.Bytes(), "job_"+jobIdStr+".obj", "binary/octet-stream", minioClient, ctx)
	if err != nil {
		return nil, err
	}

	t, err = transform(scenePath+"/"+productString+".mtl.tmpl", jobIdStr)
	if err != nil {
		return nil, err
	}
	err = upload(t.Bytes(), "job_"+jobIdStr+".mtl", "binary/octet-stream", minioClient, ctx)
	if err != nil {
		return nil, err
	}
	badgeFile, err := os.Open(scenePath + "/gcp_next_badge.png")
	if err != nil {
		return nil, err
	}
	badgeImg, err := png.Decode(badgeFile)
	if err != nil {
		return nil, err
	}

	// Add text to badge and upload to minio
	addLabel(badgeImg.(*image.NRGBA), 90, 120, req.Name)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, badgeImg)
	err = upload(buf.Bytes(), "job_"+jobIdStr+"_badge.png", "image/png", minioClient, ctx)
	if err != nil {
		return nil, err
	}

	// Add tasks to the GifJob queue for each frame to render
	var taskId int64
	for i := 0; i < 15; i++ {
		// Set up render request for each frame
		var task = renderTask{
			Frame:       int64(i),
			ProductType: req.ProductToPlug,
			Caption:     req.Name,
		}

		// Get new task id
		taskId, err = redisClient.Incr(redisContext, "counter_queued_gifjob_"+jobIdStr).Result()
		if err != nil {
			return nil, err
		}
		taskIdStr := strconv.FormatInt(taskId, 10)

		payload, err = json.Marshal(task)
		if err != nil {
			return nil, err
		}
		err = redisClient.Set(redisContext, "task_gifjob_"+jobIdStr+"_"+taskIdStr, payload, 0).Err()
		if err != nil {
			return nil, err
		}
		err = redisClient.LPush(redisContext, "gifjob_queued", jobIdStr+"_"+taskIdStr).Err()
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stdout, "enqueued gifjob_%s_%s %s\n", jobIdStr, taskIdStr, payload)
	}

	// Return job ID
	response := pb.StartJobResponse{JobId: jobIdStr}

	return &response, nil
}

func leaseNextTask() error {
	// We want to make task leasing as robust as possible. We do this by
	// shifting the task marker to a 'processing' queue that signals that we are
	// trying to work on it. Once the task is done it's removed from the
	// processing queue. If this process crashes during processing then a garbage
	// collector could move the task back into the 'queueing' queue.
	redisContext := context.Background()
	jobString, err := redisClient.BRPopLPush(redisContext, "gifjob_queued", "gifjob_processing", 0).Result()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "leased gifjob_%s\n", jobString)

	// extract task ID and job ID
	strs := strings.Split(jobString, "_")
	jobIdStr := strs[0]
	taskIdStr := strs[1]

	payload, err := redisClient.Get(redisContext, "task_gifjob_"+jobIdStr+"_"+taskIdStr).Result()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "leased gifjob_%s %s\n", jobString, payload)

	var task renderTask
	err = json.Unmarshal([]byte(payload), &task)
	if err != nil {
		return err
	}

	req := &pb.RenderRequest{
		GcsOutputBase: "jobnum" + jobIdStr, //this is the prefix for all/most objects of this job
		ObjPath:       "job_" + jobIdStr + ".obj",
		Assets: []string{
			"job_" + jobIdStr + ".mtl",
			"job_" + jobIdStr + "_badge.png",
			"k8s.png",
			"grpc.png",
		},
		Rotation:   float32(task.Frame*2 + 20),
		Iterations: 1,
	}

	_, err = renderClient.RenderFrame(context.Background(), req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error requesting frame - %v\n", err)
		return err
	}

	// delete item from gifjob_processing
	err = redisClient.LRem(redisContext, "gifjob_processing", 1, jobString).Err()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "deleted gifjob_%s\n", jobString)

	// increment "gifjob_"+jobIdStr+"_completed_counter"
	completedTaskCount, err := redisClient.Incr(redisContext, "counter_completed_gifjob_"+jobIdStr).Result()
	if err != nil {
		return err
	}
	queueLength, err := redisClient.Get(redisContext, "counter_queued_gifjob_"+jobIdStr).Result()
	if err != nil {
		return err
	}

	// if qeuedcounter = completedcounter, mark job as done
	queueLengthInt, _ := strconv.ParseInt(queueLength, 10, 64)
	fmt.Fprintf(os.Stdout, "job_gifjob_%s : %d of %d tasks done\n", jobIdStr, completedTaskCount, queueLengthInt)
	if completedTaskCount == queueLengthInt {
		finalImagePath, err := compileGifs(req.GcsOutputBase, context.Background())
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "final image path: %s\n", finalImagePath)
		var job = renderJob{
			Status:         pb.GetJobResponse_DONE,
			FinalImagePath: finalImagePath,
		}
		payloadBytes, _ := json.Marshal(job)
		err = redisClient.Set(redisContext, "job_gifjob_"+jobIdStr, payloadBytes, 0).Err()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "completed job_gifjob_%s : %d tasks\n", jobIdStr, completedTaskCount)
	}

	return nil
}

// compileGifs() will glob all minio objects prefixed with prefix, and
// stitch them together into an animated GIF, store that in minio and
// return the path of the final image
func compileGifs(prefix string, tCtx context.Context) (string, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println("error making minio client", err)
		return "", err
	}

	ctx, cancel := context.WithCancel(tCtx)

	defer cancel()

	objectCh := minioClient.ListObjects(ctx, minioBucket, minio.ListObjectsOptions{Prefix: prefix})
	var orderedObjects []minio.ObjectInfo
	for minioObj := range objectCh {
		if minioObj.Err != nil {
			return "", minioObj.Err
		}
		orderedObjects = append(orderedObjects, minioObj)
	}

	finalGif := &gif.GIF{}
	for _, frameObj := range orderedObjects {
		object, err := minioClient.GetObject(ctx, minioBucket, frameObj.Key, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		framePng, err := png.Decode(object)
		if err != nil {
			return "", err
		}

		var gifBuf bytes.Buffer
		var opt gif.Options
		opt.NumColors = 256
		err = gif.Encode(&gifBuf, framePng, &opt)

		frameGif, err := gif.Decode(&gifBuf)
		if err != nil {
			return "", err
		}

		finalGif.Image = append(finalGif.Image, frameGif.(*image.Paletted))
		finalGif.Delay = append(finalGif.Delay, 0)
	}

	finalObjName := prefix + "animated.gif"

	gifBuffer := new(bytes.Buffer)
	err = gif.EncodeAll(gifBuffer, finalGif)
	err = upload(gifBuffer.Bytes(), finalObjName, "image/gif", minioClient, ctx)
	// TODO: set final minio object to be public and return the link to it
	// instead of using a presigned URL so it's always public
	presignedURL, err := minioClient.PresignedGetObject(ctx, minioBucket, finalObjName, time.Second*24*60*60, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// Return status of job and url of image
func (server) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	var job renderJob
	statusStr, err := redisClient.Get(redisContext, "job_gifjob_"+string(req.JobId)).Result()
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stdout, "status of gifjob_%s is %s\n", req.JobId, statusStr)
	err = json.Unmarshal([]byte(statusStr), &job)
	if err != nil {
		return nil, err
	}
	response := pb.GetJobResponse{ImageUrl: job.FinalImagePath, Status: job.Status}
	return &response, nil
}

func main() {
	flag.Parse()
	port := os.Getenv("GIFCREATOR_PORT")
	if port == "" {
		if *workerMode == true {
			port = "8081"
		}
		port = "8082"
	}
	if os.Getenv("REDIS_NAME") != "" {
		redisName = os.Getenv("REDIS_NAME")
	}
	if os.Getenv("REDIS_PORT") != "" {
		redisPort = os.Getenv("REDIS_PORT")
	}
	renderName := os.Getenv("RENDER_NAME")
	if renderName == "" {
		redisName = "localhost"
	}
	renderPort := os.Getenv("RENDER_PORT")
	if renderPort == "" {
		redisPort = "8080"
	}
	renderHostAddr := renderName + ":" + renderPort
	scenePath = os.Getenv("SCENE_PATH")
	if scenePath == "" {
		scenePath = "/scene"
	}
	minioName := os.Getenv("MINIO_NAME")
	if minioName == "" {
		minioName = "localhost"
	}
	minioPort := os.Getenv("MINIO_PORT")
	if minioPort == "" {
		minioPort = "9000"
	}
	endpoint = minioName + ":" + minioPort
	accessKeyID = os.Getenv("MINIO_KEY")
	if accessKeyID == "" {
		accessKeyID = "minioaccesskeyid"
	}
	secretAccessKey = os.Getenv("MINIO_SECRET")
	if secretAccessKey == "" {
		secretAccessKey = "miniosecretaccesskey"
	}
	minioBucket = os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "gifbucket"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisName + ":" + redisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if *workerMode == true {
		// Worker mode will perpetually poll the queue and lease tasks
		fmt.Fprintf(os.Stdout, "starting gifcreator in worker mode\n")

		conn, err := grpc.Dial(renderHostAddr, grpc.WithInsecure())

		if err != nil {
			// TODO(jessup) Swap these out for proper logging
			fmt.Fprintf(os.Stderr, "cannot connect to render service %s\n%v", renderHostAddr, err)
			return
		}
		defer conn.Close()

		renderClient = pb.NewRenderClient(conn)

		for {
			err := leaseNextTask()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error working on task: %v\n", err)
			}
			time.Sleep(10 * time.Millisecond)
			// TODO(jessup) add timed sweeps for crashed jobs that never finished processing
		}
	} else {
		// Server mode will act as a gRPC server
		fmt.Fprintf(os.Stdout, "starting gifcreator in server mode\n")
		l, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("listen failed: %v", err)
		}
		srv := grpc.NewServer()
		pb.RegisterGifCreatorServer(srv, server{})
		srv.Serve(l)
	}
}
