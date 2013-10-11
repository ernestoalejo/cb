#!/bin/bash

USER=$1
PASS=$2
HOST=$3
DATE="deploy-"`date +%s`


lftp -c "set ftp:list-options -a;
set ftp:ssl-allow off;
set cmd:fail-exit yes;
open -u $USER,$PASS $HOST;
mkdir $DATE;
cd $DATE;
source temp/upload-commands
"
