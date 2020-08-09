#!/bin/bash

echo "##### Stopping API container"
docker stop cotw-api
echo "##### Copying database from container"
docker cp cotw-api:/go/src/app/leaderboard.db .
echo "##### Removing old API container"
docker rm cotw-api
echo "##### Building new image"
docker build -t cotw-api ../
echo "##### Starting new API container"
docker run -d -p 8080:8080 --name cotw-api cotw-api
echo "##### Copying database to new container"
docker cp ./leaderboard.db cotw-api:/go/src/app 
echo "##### Remove local database copy"
rm ./leaderboard.db

