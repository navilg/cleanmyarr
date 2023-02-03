#!/usr/bin/env bash

if [ "$#" == 0 ]; then
    tag=latest
else
    tag=$1
fi

docker buildx create --name dockerxbuilder --use --bootstrap
docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag linuxshots/cleanmyarr:$tag .
docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag linuxshots/cleanmyarr .