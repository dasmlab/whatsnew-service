#!/bin/bash

set -e

VERSION="0.1.0"
ARCH="amd64"

docker build -t whatsnew-service:$VERSION \
  --build-arg ARCH=$ARCH \
  --build-arg goproxy=https://proxy.golang.org \
  -f Dockerfile .

