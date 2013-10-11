#!/bin/bash
vendor/bin/phpmd app/controllers/ text ./phpmd.xml
vendor/bin/phpmd app/models text ./phpmd.xml
