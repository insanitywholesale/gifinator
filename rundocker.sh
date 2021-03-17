#!/bin/sh -x

docker network create --driver bridge gifnet
docker run --network gifnet --rm --name redis -p 6379:6379 redis:6 &
sleep 5s;
docker run --network gifnet --rm --name minio -p 9000:9000 -v /tmp/docker_data/minio/data:/data -e MINIO_ACCESS_KEY=minioaccesskeyid -e MINIO_SECRET_KEY=miniosecretaccesskey minio/minio:latest server /data &
sleep 5s;
docker run --network gifnet --rm --name render -p 8080:8080 inherently/gifinator:0.0.1 /render &
sleep 5s;
docker run --network gifnet --rm --name gifcreator-server -p 8081:8081 -e MINIO_NAME=minio -e REDIS_NAME=redis -e REDIS_PORT=6379 -e RENDER_NAME=render -e RENDER_PORT=8080 -e GIFCREATOR_PORT=8082 inherently/gifinator:0.0.1 /gifcreator &
sleep 5s;
docker run --network gifnet --rm --name gifcreator-worker -p 8082:8082 -e MINIO_NAME=minio -e REDIS_NAME=redis -e REDIS_PORT=6379 -e GIFCREATOR_PORT=8081 -e RENDER_NAME=render -e RENDER_PORT=8080 inherently/gifinator:0.0.1 /gifcreator -worker &
sleep 5s;
docker run --network gifnet --rm --name frontend -p 8090:8090 -e GIFCREATOR_NAME=gifcreator-server -e GIFCREATOR_PORT=8081 -e FRONTEND_TEMPLATES_DIR=/templates inherently/gifinator:0.0.1 /frontend &
