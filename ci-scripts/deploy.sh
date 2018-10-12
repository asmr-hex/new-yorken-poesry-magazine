#!/bin/bash
set -x # display output of commands

# deploy this thang
#
# at this point we should have the frontend bundle and prod docker-compose within
# ./prod-deploy and the docker image for the app should have been pushed to dockerhub

# turn the prod directory into a git repo and force push to deploy server
cd prod-deploy
git init
git remote add deploy "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_REPO"
git config user.name "Travis CI"
git config user.email "c@polygon.pizza"
git add .
git commit -m "deploy"
git push --force deploy master

# load updated migration files, pull latest docker image and restart on server
cmds="cd $DEPLOY_DIR;"
cmds="$cmds mv ./migrations /data/migrations;"
cmds="$cmds docker-compose pull nypm;"
cmds="$cmds docker-compose up -d"

ssh -t "$DEPLOY_USER@$DEPLOY_HOST" "$cmds" 

