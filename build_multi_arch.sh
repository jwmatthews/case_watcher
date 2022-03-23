#!/bin/bash

REPO="quay.io/jwmatthews/case_watcher"

podman manifest exists ${REPO}:latest
if [ "$?" == "0" ]
then
	podman manifest rm ${REPO}:latest
fi
podman manifest create ${REPO}:latest
podman build -t ${REPO}:amd64 --arch=amd64 -f Dockerfile .
podman build -t ${REPO}:arm64 --arch=arm64 -f Dockerfile .
podman push ${REPO}:amd64
podman push ${REPO}:arm64
podman manifest add ${REPO}:latest docker://${REPO}:amd64
podman manifest add ${REPO}:latest docker://${REPO}:arm64
podman manifest push --all ${REPO}:latest docker://${REPO}:latest