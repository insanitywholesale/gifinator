#!/bin/sh

mkdir /tmp/objcache
cp -r gifcreator/scene/ /tmp
export SCENE_PATH="/tmp/scene"
./runredis.sh &
./runminio.sh &
