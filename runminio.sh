#!/bin/sh

# 9000 is API, 9001 is web ui
docker run --rm --name minio \
	-p 9000:9000 -p 9001:9001 \
	-v /tmp/docker_data/minio/data:/data \
	-e MINIO_ROOT_USER=minioaccesskeyid \
	-e MINIO_ROOT_PASSWORD=miniosecretaccesskey \
	minio/minio:latest server --console-address ":9001" /data
