#!/bin/bash
set -x #display output of commands

# add deploy server public key to known hosts
echo "poem.computer ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBC3Ga6LBdVDCTkRvH6zH826U4Iyce4I/qjvB5RNmclEY9/eCwEjF4nNmJCgQKyyvPj/aW6yepn/n+4wveHmI7UQ=" >> "$HOME"/.ssh/known_hosts

# build frontend bundle for prod
echo "--- BUILDING FRONTEND BUNDLE FOR PROD ---"
cd client
rm -rf node_modules  # blast away old dependencies
npm install -qy      # re-install dependencies
npm run build        # this should put the bundle within ./build
cd ..

# package up bundles frontend, docker-compose.prod.yml, .env, migrations into directory
mkdir prod-deploy
mv ./client/build/* ./prod-deploy
mv docker-compose.prod.yml ./prod-deploy/docker-compose.yml
mv .env ./prod-deploy/.env
mv ./migrations ./prod-deploy/migrations

# build app docker image for prod
echo "--- BUILDING APP DOCKER IMAGE FOR PROD ---"
docker build -t "$DOCKER_REPO" .

# push app docker image to dockerhub
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USER" --password-stdin
docker push "$DOCKER_REPO"
