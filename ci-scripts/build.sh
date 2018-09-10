#!/bin/bash
set -x #display output of commands

# build frontend bundle for prod
echo "--- BUILDING FRONTEND BUNDLE FOR PROD ---"
npm install -qy
npm run build  # this should put the bundle within ./build

# build app docker image for prod
echo "--- BUILDING APP DOCKER IMAGE FOR PROD ---"
docker build -t "$DOCKER_REPO"

# push app docker image to dockerhub
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USER" --password-stdin
docker push "$DOCKER_REPO"
