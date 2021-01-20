#!/bin/sh

docker run -p 9000:9000 -v /tmp/docker_data/minio/data:/data -e MINIO_ACCESS_KEY=minioaccesskeyid -e MINIO_SECRET_KEY=miniosecretaccesskey minio/minio server /data
