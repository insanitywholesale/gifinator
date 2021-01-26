# gifinator

fork of https://github.com/GoogleCloudPlatform/gifinator to run locally and use minio for storage instead of gcs.
look at TODO for what's left to be done (hint: it's a lot).

# running (for development)
the following should be run in the order they are listed and each one in a different terminal

## redis (port 6379) and minio (port 9000)
from the root of the repo, run `./rundeps.sh`

## render (port 8080)
this one is pretty simple, use `go run render.go`

## gifcreator (worker port 8081, server port 8082)
go into its directory and run `./rungifcreator.sh`

## frontend (port 8090)
also pretty simple, from inside its directory run `./runfront.sh`

## create bucket and load assets
go to http://localhost:9000 and log in using the credentials `minioaccesskeyid` and `miniosecretaccesskey` then create a bucket named `gifbucket` and load all the assets in `gifcreator/scene` to it

## visit web page
access http://localhost:8090 using a web browser, fill in the text, select one of the 3 options, click Create and wait

# running (for production)
there are eventually going to be kubernetes manifests for it in the `infra` repository under `kube/manifests`
