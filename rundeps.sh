#!/bin/sh

mkdir /tmp/objcache
cp -r gifcreator/scene/ /tmp
./runredis.sh &
./runminio.sh &
