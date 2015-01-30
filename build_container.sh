#!/bin/bash
set -x

# build the app for linux/i386 an build the docker container
GOOS=linux GOARCH=386 go build
mkdir -p target/linux_i386
mv ./monarch target/linux_i386/
docker build -t magneticio/monarch:0.2 .
