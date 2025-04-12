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
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "gitlab.com/insanitywholesale/gifinator/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// where to find templates
	templatePath string
	//go:embed static
	staticPath embed.FS
	// gifcreator client
	gcClient pb.GifCreatorClient
	// commit info
	commitHash string
	commitDate string
)

func main() {
	templatePath = os.Getenv("FRONTEND_TEMPLATES_DIR")
	if templatePath == "" {
		templatePath = "/templates"
	}
	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		port = "8090"
	}
	gifcreatorPort := os.Getenv("GIFCREATOR_PORT")
	if gifcreatorPort == "" {
		gifcreatorPort = "8081"
	}
	gifcreatorName := os.Getenv("GIFCREATOR_NAME")
	if gifcreatorName == "" {
		gifcreatorPort = "localhost"
	}

	gcHostAddr := gifcreatorName + ":" + gifcreatorPort

	conn, err := grpc.Dial(gcHostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to gifcreator %s\n%v", gcHostAddr, err)
		return
	}
	// Store gifcreator grpc client globally
	gcClient = pb.NewGifCreatorClient(conn)
	log.Println("connected to gifcreator")

	// Create router
	m := http.NewServeMux()
	m.HandleFunc("/", handleForm)
	m.HandleFunc("/gif/", handleGif)
	m.HandleFunc("/check/", handleGifStatus)
	m.HandleFunc("/info", getInfo)
	m.Handle("/static/", http.FileServer(http.FS(staticPath)))

	// Create HTTP server with timeouts
	s := &http.Server{
		Handler:           m,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Println("about to start serving")
	// Create listener on provided port
	l, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		conn.Close()
		log.Fatal(err)
	}

	// Start HTTP server
	err = s.Serve(l)
	if err != nil {
		conn.Close()
		log.Fatal(err)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the form info, verify, and pass on
		formErrors := []string{}
		var gifName string
		var mascotType pb.Product
		err := r.ParseForm()
		if err != nil {
			// TODO Swap these out for proper logging
			fmt.Fprintf(os.Stderr, "ParseForm failure - %v", err)
		}
		if (r.Form["name"] != nil) && (len(r.Form["name"][0]) > 0) {
			gifName = r.Form["name"][0]
		} else {
			formErrors = append(formErrors, "Please provide a name")
		}
		if r.Form["mascot"] != nil {
			switch r.Form["mascot"][0] {
			case "go":
				mascotType = pb.Product_GO
			case "grpc":
				mascotType = pb.Product_GRPC
			case "kubernetes":
				mascotType = pb.Product_KUBERNETES
			default:
				mascotType = pb.Product_UNKNOWN_PRODUCT
			}
		} else {
			formErrors = append(formErrors, "Please specify a mascot")
		}
		if len(formErrors) > 0 {
			renderForm(w, formErrors)
			return
		}
		response, err := gcClient.StartJob(context.Background(), &pb.StartJobRequest{Name: gifName, ProductToPlug: mascotType})
		if err != nil {
			// TODO(jessup) Swap these out for proper logging
			fmt.Fprintf(os.Stderr, "cannot request Gif - %v", err)
			return
		}
		http.Redirect(w, r, "/gif/"+response.JobId, http.StatusMovedPermanently)
		return
	}
	renderForm(w, nil)
}

func renderForm(w http.ResponseWriter, errors []string) {
	// Show the form
	formPath := filepath.Join(templatePath, "form.html")
	layoutPath := filepath.Join(templatePath, "layout.html")

	t, err := template.ParseFiles(layoutPath, formPath)
	if err == nil {
		terr := t.ExecuteTemplate(w, "layout", errors)
		http.Error(w, terr.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type responsePageData struct {
	ImageID  string
	ImageURL string
}

func handleGif(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 2 {
		http.Error(w, "Can't find the GIF ID", http.StatusNotFound)
		return
	}

	// TODO(jessup) Look up to see if the gif has loaded. If not, show the Spinner.
	response, err := gcClient.GetJob(
		context.Background(),
		&pb.GetJobRequest{JobId: pathSegments[2]})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get status of gif - %v", err)
		return
	}

	var bodyHTMLPath string
	gifInfo := responsePageData{
		ImageID: pathSegments[2],
	}
	switch response.Status {
	case pb.GetJobResponse_PENDING:
		bodyHTMLPath = filepath.Join(templatePath, "spinner.html")
	case pb.GetJobResponse_DONE:
		bodyHTMLPath = filepath.Join(templatePath, "gif.html")
		gifInfo.ImageURL = response.ImageUrl
	default:
		bodyHTMLPath = filepath.Join(templatePath, "error.html")
	}
	layoutPath := filepath.Join(templatePath, "layout.html")

	t, err := template.ParseFiles(layoutPath, bodyHTMLPath)
	if err == nil {
		terr := t.ExecuteTemplate(w, "layout", gifInfo)
		http.Error(w, terr.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGifStatus(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 2 {
		http.Error(w, "Can't find the GIF ID", http.StatusNotFound)
		return
	}

	// TODO(jessup) Need stronger input validation here.
	response, err := gcClient.GetJob(
		context.Background(),
		&pb.GetJobRequest{JobId: pathSegments[2]})
	if err != nil {
		// TODO(jessup) Swap these out for proper logging
		fmt.Fprintf(os.Stderr, "cannot get status of gif - %v", err)
		return
	}

	jsonReponse, err := json.Marshal(response)
	if err != nil {
		// TODO(jessup) Swap these out for proper logging
		fmt.Fprintf(os.Stderr, "cannot marshal response - %v", err)
		return
	}
	w.Write(jsonReponse)
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("commitHash: " + commitHash + "\n"))
		w.Write([]byte("commitDate: " + commitDate + "\n"))
		return
	}
}
