#!/bin/bash

# Exit from script on first error
set -e

# Run composer install or update depending if it's a recurring instalation or not
if [ -f composer.lock ]; then
  php ~/bin/composer.phar update
else
  php ~/bin/composer.phar install
fi

# Update bower packages
cd client
bower install

# Remove this script
cd ..
rm post-init.sh
