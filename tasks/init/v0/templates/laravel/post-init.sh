#!/bin/bash

# Exit from script on first error
set -e

# Run composer install or update depending if it's a recurring instalation or not
php ~/bin/composer.phar install
if [ -f composer.lock ]; then
  php ~/bin/composer.phar update
fi

# Update bower packages
cd client
bower install

# Remove this script
cd ..
rm post-init.sh
