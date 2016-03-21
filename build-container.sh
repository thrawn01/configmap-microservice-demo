#!/bin/sh

DOCKER_REPO=thrawn01
docker build -t ${DOCKER_REPO}/configmap-microservice-demo:latest .
docker push ${DOCKER_REPO}/configmap-microservice-demo:latest
