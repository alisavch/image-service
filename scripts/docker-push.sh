#!/bin/bash

set -e -x

docker build -t alisavch/api:latest -f Dockerfile .
docker build -t alisavch/consumer:latest -f Dockerfile-consumer .

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push alisavch/api:latest
docker push alisavch/consumer:latest