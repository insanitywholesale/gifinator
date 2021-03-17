#!/bin/sh -x

docker run --rm --name minio -p 9000:9000 -v /tmp/docker_data/minio/data:/data -e MINIO_ACCESS_KEY=minioaccesskeyid -e MINIO_SECRET_KEY=miniosecretaccesskey minio/minio:latest server /data &
docker run --rm --name redis -p 6379:6379 redis:6 &
docker run --rm --name render -p 8080:8080 inherently/gifinator:0.0.1 /render &
docker run --rm --name gifcreator-server -p 8081:8081 -e REDIS_NAME=redis -e REDIS_PORT=6379 -e GIFCREATOR_PORT=8081 inherently/gifinator:0.0.1 /gifcreator &
docker run --rm --name gifcreator-worker -p 8082:8082 -e REDIS_NAME=redis -e REDIS_PORT=6379 -e GIFCREATOR_PORT=8082 -e RENDER_NAME=render -e RENDER_PORT=8080 inherently/gifinator:0.0.1 /gifcreator -worker &
docker run --rm --name frontend -p 8090:8090 -e GIFCREATOR_NAME=gifcreator-server -e GIFCREATOR_PORT=8082 inherently/gifinator:0.0.1 /frontend &
