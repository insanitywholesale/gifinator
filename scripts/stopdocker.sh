#!/bin/sh

#stop containers
docker stop \
	frontend \
	gifcreator-worker \
	gifcreator-server \
	render \
	minio \
	redis

#remove containers
docker rm \
	frontend \
	gifcreator-worker \
	gifcreator-server \
	render \
	minio \
	redis
