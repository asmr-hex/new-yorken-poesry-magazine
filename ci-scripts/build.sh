#!/bin/bash
set -x #display output of commands

# add deploy ip to known hosts
ssh-keyscan -t "$TRAVIS_SSH_KEY_TYPES" -H "$DEPLOY_IP_ADDRESS" 2>&1 | tee -a "$HOME"/.ssh/known_hosts

# build frontend bundle for prod
echo "--- BUILDING FRONTEND BUNDLE FOR PROD ---"
cd client
npm install -qy
npm run build  # this should put the bundle within ./build
cd ..

# package up bundles frontend and docker-compose.prod.yml into directory
mkdir prod-deploy
mv ./client/build/* ./prod-deploy
mv docker-compose.prod.yml ./prod-deploy/docker-compose.yml

# build app docker image for prod
echo "--- BUILDING APP DOCKER IMAGE FOR PROD ---"
docker build -t "$DOCKER_REPO" .

# push app docker image to dockerhub
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USER" --password-stdin
docker push "$DOCKER_REPO"
