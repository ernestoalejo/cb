#!/bin/bash

USER=$1
PASS=$2
HOST=$3


cd temp
rm hashes


lftp -c "set ftp:list-options -a;
set ftp:ssl-allow off;
set cmd:fail-exit yes;
open -u $USER,$PASS $HOST;
find push-hashes || exit 101;
get push-hashes -o hashes;
" 2> /dev/null || true
