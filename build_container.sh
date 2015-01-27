#!/bin/bash
set -x

# build the app for linux/i386 an build the docker container
GOOS=linux GOARCH=386 go build
mv ./mesos-tester target/linux_i386/
docker build -t magneticio/tester:6.1 .