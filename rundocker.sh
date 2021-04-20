#!/bin/sh -x

IMG="inherently/gifinator:0.0.4"
NET="gifnet"

docker pull $IMG
sleep 5s
docker pull redis:6
sleep 5s
docker pull minio:latest
docker network create --driver bridge $NET
docker run --network $NET --rm --name redis -p 6379:6379 redis:6 &
sleep 5s
docker run --network $NET --rm --name minio \
	-p 9000:9000 \
	-v /tmp/docker_data/minio/data:/data \
	-e MINIO_ACCESS_KEY=minioaccesskeyid \
	-e MINIO_SECRET_KEY=miniosecretaccesskey \
	minio/minio:latest server /data &
sleep 5s
docker run --network $NET --rm --name render \
	-p 8080:8080 \
	-e MINIO_NAME=minio $IMG /render &
sleep 5s
docker run --network $NET --rm --name gifcreator-server \
	-p 8082:8082 \
	-e MINIO_NAME=minio \
	-e REDIS_NAME=redis \
	-e REDIS_PORT=6379 \
	-e RENDER_NAME=render \
	-e RENDER_PORT=8080 $IMG /gifcreator &
sleep 5s
docker run --network $NET --rm --name gifcreator-worker \
	-p 8081:8081 \
	-e MINIO_NAME=minio \
	-e REDIS_NAME=redis \
	-e REDIS_PORT=6379 \
	-e RENDER_NAME=render \
	-e RENDER_PORT=8080 $IMG /gifcreator -worker &
sleep 5s
docker run --network $NET --rm --name frontend \
	-p 8090:8090 \
	-e FRONTEND_PORT=8090 \
	-e GIFCREATOR_NAME=gifcreator-server \
	-e GIFCREATOR_PORT=8082 \
	-e FRONTEND_TEMPLATES_DIR=/templates $IMG /frontend &
