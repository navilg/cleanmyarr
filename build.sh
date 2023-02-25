#!/usr/bin/env bash

if [ "$#" == 0 ]; then
    echo "Usage:"
    echo "bash $0 tag"
    echo "For e.g."
    echo "bash $0 1.2.3"
    exit 1
fi

tag=$1

docker buildx create --name dockerxbuilder --use --bootstrap
docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag linuxshots/cleanmyarr:$tag .
docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag linuxshots/cleanmyarr .