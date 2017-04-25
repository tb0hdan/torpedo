#!/bin/bash

# Get yours at https://api.slack.com/custom-integrations/legacy-tokens
# and store as
# TOKEN="xxx"
# inside token.sh
source ./token.sh
go run *.go -token ${TOKEN} -lastfm_key ${LASTFM_KEY} -lastfm_secret ${LASTFM_SECRET}

