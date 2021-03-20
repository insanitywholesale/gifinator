# gifinator

fork of https://github.com/GoogleCloudPlatform/gifinator to run locally and use minio for storage instead of gcs.
look at TODO for what's left to be done.

# running (for development)
the following should be run in the order they are listed and each one in a different terminal

## docker
### configure name resolution
add the following in `/etc/hosts`
```
127.0.0.1   minio
127.0.1.1   minio
```

### start containers
run `./rundocker.sh` from the root of the repo to start all the required containers

### create bucket and load assets
go to `http://localhost:9000` and log in using the credentials `minioaccesskeyid` and `miniosecretaccesskey` then create a bucket named `gifbucket` and load all the assets in `gifcreator/scene` to it

### visit web frontend
access http://localhost:8090 using a web browser, fill in the text, select one of the 3 options, click `Create` and wait

## manually
### redis (port 6379) and minio (port 9000)
from the root of the repo, run `./rundeps.sh`

### render (port 8080)
this one is pretty simple, use `go run render.go`

### gifcreator (worker port 8081, server port 8082)
go into its directory and run `./rungifcreator.sh`

### frontend (port 8090)
also pretty simple, from inside its directory run `./runfront.sh`

### create bucket and load assets
go to `http://localhost:9000` and log in using the credentials `minioaccesskeyid` and `miniosecretaccesskey` then create a bucket named `gifbucket` and load all the assets in `gifcreator/scene` to it

## visit web page
access http://localhost:8090 using a web browser, fill in the text, select one of the 3 options, click `Create` and wait

# running (for production)
there are kubernetes manifests for it in [here in my infra repository](https://gitlab.com/insanitywholesale/infra/-/tree/master/kube/manifests/gifinator)

# documentation
env vars for each service are listed below
## redis
see [official docs](https://github.com/librenms/docker/blob/263c47e895850e6c7a4cafedd73fadd43b870711/doc/docker/environment-variables.md)
## minio
see [official docs](https://github.com/minio/minio/tree/9171d6ef651a852b48f39f828c3d01e30fbf4e9c/docs/config)
## render
- `RENDER_PORT`: port the service will run at
- `MINIO_NAME`: minio server domain name
- `MINIO_PORT`: minio server port number
- `MINIO_BUCKET`: minio bucket to be used
- `MINIO_KEY`: minio access key
- `MINIO_SECRET`: minio secret key
## gifcreator (server mode)
- `GIFCREATOR_PORT`: port the service will run at
- `SCENE_PATH`: path to find files for gif creation
- `MINIO_NAME`: minio server domain name
- `MINIO_PORT`: minio server port number
- `MINIO_BUCKET`: minio bucket to be used
- `MINIO_KEY`: minio access key
- `MINIO_SECRET`: minio secret key
- `REDIS_NAME`: redis server domain name
- `REDIS_PORT`: redis server port number
- `RENDER_NAME`: renderer domain name
- `RENDER_PORT`: renderer port number
## gifcreator (worker mode)
- `GIFCREATOR_PORT`: port the service will run at
- `SCENE_PATH`: path to find files for gif creation
- `MINIO_NAME`: minio server domain name
- `MINIO_PORT`: minio server port number
- `MINIO_BUCKET`: minio bucket to be used
- `MINIO_KEY`: minio access key
- `MINIO_SECRET`: minio secret key
- `REDIS_NAME`: redis server domain name
- `REDIS_PORT`: redis server port number
- `RENDER_NAME`: renderer domain name
- `RENDER_PORT`: renderer port number
## frontend
- `FRONTEND_PORT`: port the service will run at
- `FRONTEND_TEMPLATES_DIR`: directory for html templates
- `GIFCREATOR_NAME`: gifcreator server domain name
- `GIFCREATOR_PORT`: gifcreator server port number

# legal stuff
the original is [here](https://github.com/GoogleCloudPlatform/gifinator) and its license and legal stuff apply, I'm not trying to steal anything
