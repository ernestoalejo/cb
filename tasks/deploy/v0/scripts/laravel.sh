#!/bin/bash

# Script args
basename=$1

# Prepare deploy folder
rm -rf temp/deploy
mkdir temp/deploy

# Copy client dist files
cp -r dist temp/deploy/public_html

# Then copy public_html folder
cp -r ../public_html temp/deploy

# Move the app folder excluding the storage folder
rsync -aq --exclude=app/storage/ ../app temp/deploy

# Move the rest of Laravel folders
cp -r ../bootstrap temp/deploy
cp -r ../vendor temp/deploy

# Replace the dev templates with the processed ones
rm -r temp/deploy/app/views
mv temp/laravel-templates temp/deploy/app/views

# Base template has a special processing
rm temp/deploy/app/views/$basename
mv temp/$basename temp/deploy/app/views/$basename

# Move deploy to the root folder
rm -rf ../deploy
mv temp/deploy ..
