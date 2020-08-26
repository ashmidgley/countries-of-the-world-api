#!/bin/bash

echo "##### Stopping API container"
docker stop cotw-api
echo "##### Removing old API container"
docker rm cotw-api
echo "##### Building new image"
docker build -t cotw-api ../
echo "##### Starting new API container"
docker run -d -p 8080:8080 -v countries-of-the-world:/go/src/app --name cotw-api cotw-api
