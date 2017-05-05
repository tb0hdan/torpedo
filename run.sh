#!/bin/bash

# Get yours at https://api.slack.com/custom-integrations/legacy-tokens
# and store as
# TOKEN="xxx"
# inside token.sh
if [ "$1" != "" ]; then
	FNAME="./token.sh.$1"
else
	FNAME="./token.sh"
fi

source ${FNAME}

GOPATH=$(pwd) go run src/torpedobot/*.go -token ${TOKEN} -telegram ${TELEGRAM} -lastfm_key ${LASTFM_KEY} -lastfm_secret ${LASTFM_SECRET} -google_webapp_key ${GOOGLE_WEBAPP_KEY} -pinterest_token ${PINTEREST_TOKEN}
