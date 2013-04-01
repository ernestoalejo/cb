# Makefile to simplify some common AppEngine actions.
# Use 'make help' for a list of commands.

# Email & Password for the AppEngine account will be read
# from this files.
EMAIL = `cat ~/.email`
PASSWORD = `cat ~/.psw`

# Path that will not show up in the console when requested
EXCLUDE_PATHS = -e 'GET /static' -e 'GET /favicon.ico' -e 'GET /_ah/img' -e '/_ah/stats'

# The file containing a copy of the production data.
# See the backup action.
BACKUP = ~/Backup/
CODE_BACKUP = ~/CodeBackup/

# The name of the project, and it's base path.
PROJECT = `cat app.yaml | sed -n 's/^application: *//p'`
PROJECT_VERSION = `cat app.yaml | sed -n 's/^version: *//p'`
PROJECT_PATH = `pwd`

# The Python binary and the AppEngine script paths.
PYTHON = python
GAE_PATH = ~/google_appengine
APPCFG = @echo $(PASSWORD) | $(PYTHON) $(GAE_PATH)/appcfg.py --email=$(EMAIL) --passin
DEV_APPSERVER = @$(PYTHON) $(GAE_PATH)/dev_appserver.py --datastore_path=datastore

# The URLs to access the remote app.
REMOTE_API_URL = http://python.$(PROJECT).appspot.com/_ah/remote_api
LOCAL_REMOTE_API_URL = http://localhost:8080/_ah/remote_api

# The temp directory.
TMP = /tmp

# The configuration file
CONF_FILE=conf/conf.go

RM=rm -rf

# The default command.
default: help


# Show a list of commands.
help:
	@echo "Available commands:"
	@sed -n '/^[a-zA-Z0-9_.]*:/s/:.*//p' <Makefile | sort


# Serve the app in localhost:8080.
serve:
	$(DEV_APPSERVER) ./ 2>&1 | grep --line-buffered -v $(EXCLUDE_PATHS) 2>&1 | egrep --color=always '^|WARNING.*|ERROR.*'


# Alias for the serve target
server: serve


# Serve the app in localhost:8080. When an email is sent,
# it's logged to the console.
serve_mail:
	$(DEV_APPSERVER) --show_mail_body ./


# Serve the app in localhost:8080, allowing for external connections
# to see and interact with it.
serve_remote:
	$(DEV_APPSERVER) --address 0.0.0.0 ./


# Removes the local datastore file.
clear_datastore:
	$(RM) ./datastore


# Deploy the app to AppEngine.
deploy:
	@echo "Updating $(PROJECT)"
	$(APPCFG) $(APPCFG_FLAGS) update ./


# Rollback a failed deploy.
rollback:
	$(APPCFG) rollback ./


# Deploy only the indexes information.
update_indexes:
	$(APPCFG) update_indexes ./


# Deploy only the cron information.
update_cron:
	$(APPCFG) update_cron ./


# Deploy only the taskqueues information.
update_queues:
	$(APPCFG) update_queues ./


# Starts a cleaning process over the non-used indexes.
vacuum_indexes:
	$(APPCFG) vacuum_indexes ./


# Download the production data and send it to the local app
# (you have to start it before in other terminal).
sync_datastore:
	$(RM) /tmp/gae
	$(APPCFG) download_data --url=$(REMOTE_API_URL) --file=$(TMP)/gae ./ -A "s~$(PROJECT)"
	$(APPCFG) upload_data --url=$(LOCAL_REMOTE_API_URL) --file=$(TMP)/gae ./ -A "dev~$(PROJECT)"
	$(RM) bulkloader*


# Download the production data and send a copy to the desktop, to
# backup all the data.
backup:
	mkdir -p $(BACKUP)

	$(APPCFG) download_data --url=$(REMOTE_API_URL) --file=$(TMP)/gae ./ -A "s~$(PROJECT)" --kind=User
	$(RM) bulkloader*
	@mv $(TMP)/gae $(BACKUP)/User.backup

	$(APPCFG) download_data --url=$(REMOTE_API_URL) --file=$(TMP)/gae ./ -A "s~$(PROJECT)" --kind=Film
	$(RM) bulkloader*
	@mv $(TMP)/gae $(BACKUP)/Film.backup


# Format the Go code.
format:
	gofmt -w .


# Install all the necessary packages to test the server version of
# the app.
prepare_test:
	@go test -i ./modules


# Test the server part of the app.
test:
	@go test ./modules	-test.v


# Download the source code of the app from the cloud.
download:
	mkdir -p $(CODE_BACKUP)
	$(APPCFG) download_app -A $(PROJECT) -V $(PROJECT_VERSION) $(CODE_BACKUP)
 

# To decrypt the configurations
decrypt_conf:
	openssl cast5-cbc -d -in ${CONF_FILE}.cast5 -out ${CONF_FILE}
	chmod 600 ${CONF_FILE}
 

# To encrypt the configurations
encrypt_conf:
	openssl cast5-cbc -e -in ${CONF_FILE} -out ${CONF_FILE}.cast5


