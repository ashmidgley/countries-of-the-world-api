#!/bin/bash

date=$(date '+%Y-%m-%d')
container='cotw-api'
database='leaderboard'
space=''
dir='/root/scripts'

echo "##### Copying backup file from container to local directory"
docker cp $container:/go/src/app/$database.db $dir

echo "##### Zipping up db file"
zip $dir/$date-leaderboard.zip $dir/$database.db

echo "##### Moving zipped backup to Digital Ocean space"
s3cmd put $dir/$date-leaderboard.zip s3://$space

echo "##### Deleting local backup files"
rm $dir/$database.db
rm $dir/$date-leaderboard.zip

