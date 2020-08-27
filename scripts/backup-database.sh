#!/bin/bash

date=$(date '+%Y-%m-%d')
database='leaderboard'
dir=''
space=''

echo "##### Zipping up db file"
zip $dir/$date-leaderboard.zip $dir/$database.db

echo "##### Moving zipped backup to Digital Ocean space"
s3cmd put $dir/$date-leaderboard.zip s3://$space

echo "##### Deleting local backup file"
rm $dir/$date-leaderboard.zip
