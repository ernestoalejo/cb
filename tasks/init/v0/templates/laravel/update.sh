#!/bin/bash

if [ "$1" == "" ]; then
  echo "first arg should be the user name"
  exit 1
fi

USER=$1
HOST="ftp://example.com"

read -s -p "Enter Password: " PASS
echo ""

if [ "$2" == "init" ]; then
  git ftp init --user $1 --passwd $PASS $HOST --syncroot deploy
else
  git ftp push --user $1 --passwd $PASS $HOST --syncroot deploy
fi
