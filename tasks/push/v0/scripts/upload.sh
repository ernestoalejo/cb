#!/bin/bash

USER=$1
PASS=$2
HOST=$3


lftp -c "set ftp:list-options -a;
set ftp:ssl-allow off;
set cmd:fail-exit yes;
open -u $USER,$PASS $HOST;
cd example || (mkdir example && cd example);
source temp/upload-commands
"
