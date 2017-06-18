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

while :; do
    make build_only
    bin/torpedobot -slack ${SLACK} -telegram ${TELEGRAM} -lastfm_key ${LASTFM_KEY} -lastfm_secret ${LASTFM_SECRET} -google_webapp_key ${GOOGLE_WEBAPP_KEY} -pinterest_token ${PINTEREST} -jabber ${JABBER} -skype ${SKYPE} -facebook ${FACEBOOK} -kik ${KIK} -kik_webhook_url ${KIK_WEBHOOK_URL} -line ${LINE} -matrix ${MATRIX} -soundcloud_id ${SOUNDCLOUD_ID}
    sleep 3
done
