#!/bin/bash
set -x #display output of commands

# add deploy server public key to known hosts
echo "poem.computer ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBC3Ga6LBdVDCTkRvH6zH826U4Iyce4I/qjvB5RNmclEY9/eCwEjF4nNmJCgQKyyvPj/aW6yepn/n+4wveHmI7UQ=" >> "$HOME"/.ssh/known_hosts

# build frontend bundle for prod
echo "--- BUILDING FRONTEND BUNDLE FOR PROD ---"
# for some reason, travis ci won't let us build within the original client
# dir, so we need to,
# mkdir prod-client             # make a new prod-client build directory
# cd prod-client
# cp ../client/package*.json .  # only copy over the package*.json
# ls
# npm install                   # install all dependencies
# ls
# cp ../client/* .              # after installation, copy over src
cd client  # debugging
npm run build                 # this should put the bundle within ./build
cd ..

# package up bundles frontend, docker-compose.prod.yml, .env into directory
mkdir prod-deploy
mv ./client/build/* ./prod-deploy
mv docker-compose.prod.yml ./prod-deploy/docker-compose.yml
mv .env ./prod-deploy/.env

# build app docker image for prod
echo "--- BUILDING APP DOCKER IMAGE FOR PROD ---"
docker build -t "$DOCKER_REPO" .

# push app docker image to dockerhub
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USER" --password-stdin
docker push "$DOCKER_REPO"
