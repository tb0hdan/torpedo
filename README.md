# torpedo
Slack bot

# Intro

Some of the code is loosely based on github.com/nlopes/slack


# Installation

```
go get -u github.com/nlopes/slack
go get -u github.com/shkh/lastfm-go/lastfm
go get -u golang.org/x/net/html
go get -u golang.org/x/text/encoding/charmap
go get -u golang.org/x/text/transform
```

Get token:

`https://api.slack.com/custom-integrations/legacy-tokens`

Paste token as `gotestbot/token.sh`

```bash
TOKEN="xxxttt"
LASTFM_KEY="aaa"
LASTFM_SECRET="bbb"
```

# Running

```bash
cd gotestbot
./run.sh
```
