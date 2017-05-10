# torpedo
Jabber/Skype/Slack/Telegram bot

# Intro

Torpedo uses multiple Jabber/Skype/Slack/Telegram accounts (at least one is required)


# Configuration

Get Slack token(s):

`https://api.slack.com/custom-integrations/legacy-tokens`

Paste token as `token.sh`

Get Telegram/Jabber accounts.

Get Skype channel creds (https://dev.botframework.com/)

Optional parameters (all or any combination of)

```bash
TOKEN="xxxttt,aaabbb"
TELEGRAM="xxx,yyy"
JABBER="user@host.com:supersecret,user2@anotherhost.com:a1FvH12"
SKYPE="app_id:app_password,app_id2:app_password2"
```


Mandatory parameters:


```bash
LASTFM_KEY="aaa"
LASTFM_SECRET="bbb"
GOOGLE_WEBAPP_KEY="ccc"
PINTEREST_TOKEN="ddd"
```

# Running

```bash
make deps
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
