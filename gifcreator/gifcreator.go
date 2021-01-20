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
	"github.com/golang/freetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"golang.org/x/image/font/gofont/gobold"
	"google.golang.org/grpc"
	"gopkg.in/redis.v5" //very outdated api version
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

const serviceName = "gifcreator"

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
	redisClient   *redis.Client
	renderClient  pb.RenderClient
	scenePath     string
	deploymentId  string
	workerMode    = flag.Bool("worker", false, "run in worker mode rather than server")
	gcsBucketName string //TODO: minio
)

// inputPath kinda sus
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

func upload(outBytes []byte, outputPath string, mimeType string, client *minio.Client, ctx context.Context) error {
	log.Println("outputPath:", outputPath)
	//objName := strings.TrimLeft(outputPath[strings.LastIndex(outputPath, "/"):], "/")
	objName := outputPath
	uploadInfo, err := client.PutObject(ctx, "gifbucket", objName, bytes.NewReader(outBytes), int64(len(outBytes)), minio.PutObjectOptions{ContentType: mimeType})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return nil
}

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
	// Retrieive the next job ID from Redis
	jobId, err := redisClient.Incr("gifjob_counter").Result()
	if err != nil {
		return nil, err
	}
	jobIdStr := strconv.FormatInt(jobId, 10)

	// Create a new RenderJob queue for that job
	var job = renderJob{
		Status: pb.GetJobResponse_PENDING,
	}
	payload, _ := json.Marshal(job)
	err = redisClient.Set("job_gifjob_"+jobIdStr, payload, 0).Err()
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New("truenas.hell:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("katie", "Asus_hol1", ""),
		Secure: false,
	})
	if err != nil {
		log.Println("error making minio client", err)
		return nil, err
	}

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

	// Generate the assets needed to render the frame, and push them to GCS
	// TODO: this entire section kinda sus

	t, err := transform(scenePath+"/"+productString+".obj.tmpl", jobIdStr)
	if err != nil {
		return nil, err
	}
	err = upload(t.Bytes() /*this is outputPath*/, "job_"+jobIdStr+".obj", "binary/octet-stream", minioClient, ctx)
	if err != nil {
		return nil, err
	}
	t, err = transform(scenePath+"/"+productString+".mtl.tmpl", jobIdStr)
	if err != nil {
		return nil, err
	}
	err = upload(t.Bytes() /*this is outputPath*/, "job_"+jobIdStr+".mtl", "binary/octet-stream", minioClient, ctx)
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
	addLabel(badgeImg.(*image.NRGBA), 90, 120, req.Name)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, badgeImg)
	err = upload(buf.Bytes() /*this is outputPath*/, "job_"+jobIdStr+"_badge.png", "image/png", minioClient, ctx)
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

		//Get new task id
		taskId, err = redisClient.Incr("counter_queued_gifjob_" + jobIdStr).Result()
		if err != nil {
			return nil, err
		}
		taskIdStr := strconv.FormatInt(taskId, 10)

		payload, err = json.Marshal(task)
		if err != nil {
			return nil, err
		}
		err = redisClient.Set("task_gifjob_"+jobIdStr+"_"+taskIdStr, payload, 0).Err()
		if err != nil {
			return nil, err
		}
		err = redisClient.LPush("gifjob_queued", jobIdStr+"_"+taskIdStr).Err()
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
	/**
	 * We want to make task leasing as robust as possible. We do this by
	 * shifting the task marker to a 'processing' queue that signals that we are
	 * trying to work on it. Once the task is done it's removed from the
	 * processing queue. If this process crashes during processing then a garbage
	 * collector could move the task back into the 'queueing' queue.
	 */
	jobString, err := redisClient.BRPopLPush("gifjob_queued", "gifjob_processing", 0).Result()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "leased gifjob_%s\n", jobString)

	// extract task ID and job ID
	strs := strings.Split(jobString, "_")
	jobIdStr := strs[0]
	taskIdStr := strs[1]

	payload, err := redisClient.Get("task_gifjob_" + jobIdStr + "_" + taskIdStr).Result()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "leased gifjob_%s %s\n", jobString, payload)

	var task renderTask
	err = json.Unmarshal([]byte(payload), &task)
	if err != nil {
		return err
	}

	// TODO: path stuff and render request with gcsBucketName absolutely offenders
	outputPrefix := "out." + jobIdStr
	//outputBasePath := gcsBucketName + "/" + outputPrefix
	req := &pb.RenderRequest{
		GcsOutputBase: "gifbucket",
		ObjPath:       "job_" + jobIdStr + ".obj",
		Assets: []string{
			"job_" + jobIdStr + ".mtl",
			"job_" + jobIdStr + "_badge.png",
			"k8s.png",
			"grpc.png",
			/*
				gcsBucketName + "/job_" + jobIdStr + ".mtl",
				gcsBucketName + "/job_" + jobIdStr + "_badge.png",
				gcsBucketName + "/k8s.png",
				gcsBucketName + "/grpc.png",
			*/
		},
		Rotation:   float32(task.Frame*2 + 20),
		Iterations: 1,
	}

	_, err = renderClient.RenderFrame(context.Background(), req)

	if err != nil {
		// TODO(jessup) Swap these out for proper logging
		fmt.Fprintf(os.Stderr, "error requesting frame - %v\n", err)
		return err
	}

	// delete item from gifjob_processing
	err = redisClient.LRem("gifjob_processing", 1, jobString).Err()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "deleted gifjob_%s\n", jobString)

	// increment "gifjob_"+jobIdStr+"_completed_counter"
	completedTaskCount, err := redisClient.Incr("counter_completed_gifjob_" + jobIdStr).Result()
	if err != nil {
		return err
	}
	queueLength, err := redisClient.Get("counter_queued_gifjob_" + jobIdStr).Result()
	if err != nil {
		return err
	}
	// if qeuedcounter = completedcounter, mark job as done
	queueLengthInt, _ := strconv.ParseInt(queueLength, 10, 64)
	fmt.Fprintf(os.Stdout, "job_gifjob_%s : %d of %d tasks done\n", jobIdStr, completedTaskCount, queueLengthInt)
	if completedTaskCount == queueLengthInt {
		finalImagePath, err := compileGifs(outputPrefix, context.Background())
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "final image path: %s\n", finalImagePath)
		var job = renderJob{
			Status:         pb.GetJobResponse_DONE,
			FinalImagePath: finalImagePath,
		}
		payloadBytes, _ := json.Marshal(job)
		err = redisClient.Set("job_gifjob_"+jobIdStr, payloadBytes, 0).Err()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "completed job_gifjob_%s : %d tasks\n", jobIdStr, completedTaskCount)
	}

	return nil
}

/**
 * compileGifs() will glob all GCS objects prefixed with prefix, and
 * stitch them together into an animated GIF, store that in GCS and return the
 * path of the final image
 */
func compileGifs(prefix string, tCtx context.Context) (string, error) {
	minioClient, err := minio.New("truenas.hell:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("katie", "Asus_hol1", ""),
		Secure: false,
	})
	if err != nil {
		log.Println("error making minio client", err)
		return "", err
	}

	ctx, cancel := context.WithCancel(tCtx)

	defer cancel()

	objectCh := minioClient.ListObjects(ctx, "gifbucket", minio.ListObjectsOptions{Prefix: prefix})

	var orderedObjects []minio.ObjectInfo
	for minioObj := range objectCh {
		if minioObj.Err != nil {
			fmt.Println("object error:", minioObj.Err)
		}
		orderedObjects = append(orderedObjects, minioObj) //might need pointer to object instead of just object
	}

	/*
		it := gcsClient.Bucket(gcsBucketName).Objects(tCtx, &storage.Query{Prefix: prefix, Versions: false})
		// Results from GCS are unordered, so pull the list into memory and sort it
		var orderedObjects []storage.ObjectAttrs
		for {
			obj, err := it.Next()
			if err == iterator.Done {
				break
			}
			orderedObjects = append(orderedObjects, *obj)
		}
		slice.Sort(orderedObjects[:], func(i, j int) bool {
			return orderedObjects[i].Name < orderedObjects[j].Name
		})
	*/

	finalGif := &gif.GIF{}
	for _, frameObj := range orderedObjects {
		object, err := minioClient.GetObject(ctx, "gifbucket", frameObj.Key, minio.GetObjectOptions{})
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

	finalObjName := prefix + "/animated.gif"
	//finalGifBytes := bytes.NewReader(finalGif)
	//finalGifBytesIOReader := bytes.NewReader([]byte(finalGif))
	// TODO: need to save the gif struct as object
	//uploadInfo, err := minioClient.PutObject(ctx, "gifbucket", finalObjName, /*need ioreader and size for gif*/, minio.PutObjectOptions{ContentType: "image/gif"})
	//if err != nil {
	//	fmt.Println(err)
	//	return "", err
	//}
	//fmt.Println("Successfully uploaded bytes: ", uploadInfo)

	// TODO: set final minio object to be public and return the link to it
	//return gcsBucketName + ".storage.googleapis.com/" + finalObjName, nil
	// TODO: don't ignore the above TODO and change what is returned
	return "http://truenas.hell:9000/minio/gifbucket/" + finalObjName, nil
}

func (server) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	var job renderJob
	statusStr, err := redisClient.Get("job_gifjob_" + string(req.JobId)).Result()
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
	port = "9191"
	i, err := strconv.Atoi(port)
	if (err != nil) || (i < 1) {
		log.Fatalf("please set env var GIFCREATOR_PORT to a valid port")
	}

	// TODO(jessup) Need stricter checking here.
	redisName := os.Getenv("REDIS_NAME")
	redisPort := os.Getenv("REDIS_PORT")
	//projectID := os.Getenv("GOOGLE_PROJECT_ID")
	renderName := os.Getenv("RENDER_NAME")
	renderPort := os.Getenv("RENDER_PORT")
	renderHostAddr := renderName + ":" + renderPort
	deploymentId = os.Getenv("DEPLOYMENT_ID")
	gcsBucketName = os.Getenv("GCS_BUCKET_NAME")
	scenePath = os.Getenv("SCENE_PATH")

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
