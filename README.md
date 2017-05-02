# torpedo
Slack bot

# Intro

Torpedo uses multiple Slack accounts (at least one is required)

Some of the code is loosely based on github.com/nlopes/slack


# Installation

```
go get -u gopkg.in/h2non/filetype.v1
go get -u github.com/nlopes/slack
go get -u github.com/shkh/lastfm-go/lastfm
go get -u golang.org/x/net/html
go get -u golang.org/x/text/encoding/charmap
go get -u golang.org/x/text/transform
go get -u github.com/google/google-api-go-client/googleapi/transport
go get -u google.golang.org/api/youtube/v3
```

Get Slack token(s):

`https://api.slack.com/custom-integrations/legacy-tokens`

Paste token as `gotestbot/token.sh`

```bash
TOKEN="xxxttt,aaabbb"
LASTFM_KEY="aaa"
LASTFM_SECRET="bbb"
GOOGLE_WEBAPP_KEY="ccc"
PINTEREST_TOKEN="ddd"
```

# Running

```bash
cd gotestbot
./run.sh
```

# Commands

Help:

```
!?
!h
!help
```

Encoding and crypto:

`!b64e`   - Base64 encode

`!b64d`   - Base64 decode

`!md5`    - MD5 hash

`!sha1`   - SHA1 hash

`!sha256` - SHA256 hash

`!sha512` - SHA512 hash

Multimedia:

`!lastfm` - Last.FM artist/tag search

`!youtube` - Search Youtube, Track name -> Video URL

`!bashim` - Bash.Im random quote

`!bashorg` - Bash.Org random quote

`!qr` - String to QR using Google API

`!wiki` - Wikipedia search

`!pinterest` - Pinterest boards

`!tinyurl` - TinyURL shortener
