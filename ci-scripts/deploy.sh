#!/bin/bash
set -x # display output of commands

# deploy this thang
#
# at this point we should have the frontend bundle built within
# ./build and the docker image for the app should have been pushed
# to docker hub

# create a new deploy repo
mkdir prod-deploy
mv ./build/* ./prod-deploy
mv docker-compose.prod.yml ./prod-deploy/docker-compose.yml
cd prod-deploy
git init
git remote add deploy "$DEPLOY_HOST:DEPLOY_APP_DIR"
git config user.name "Travis CI"
git config user.email "c@polygon.pizza"
git add .
git commit -m "deploy"
git push --force deploy master

# pull latest docker image and restart on server
ssh -t "$DEPLOY_HOST" "cd $DEPLOY_APP_DIR; docker-compose pull nypm; docker-compose up -d" 

