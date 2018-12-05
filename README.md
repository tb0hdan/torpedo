torpedo - Pluggable, multi-network asynchronous chat bot written in Go
======

[![Build Status](https://api.travis-ci.org/tb0hdan/torpedo.svg?branch=master)](https://travis-ci.org/tb0hdan/torpedo)
[![Go Report Card](https://goreportcard.com/badge/github.com/tb0hdan/torpedo)](https://goreportcard.com/report/github.com/tb0hdan/torpedo)
[![codecov](https://codecov.io/gh/tb0hdan/torpedo/branch/master/graph/badge.svg)](https://codecov.io/gh/tb0hdan/torpedo)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftb0hdan%2Ftorpedo.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftb0hdan%2Ftorpedo?ref=badge_shield)


# Intro

Torpedo uses multiple accounts (at least one is required). Supported transports:

- Facebook
- Jabber
- Kik
- Line
- Matrix (matrix.org only atm)
- Skype (via BotFramework)
- Microsoft Teams (via BotFramework and CustomBots webhook)
- Slack
- Telegram
- IRC

# See it in action

Jabber: torpedobot@jabber.cz

Skype: https://join.skype.com/bot/f61c6815-438d-4795-8aaa-9b1d8d2a342a

Telegram: http://t.me/TorpedoTelegramBot

Line:

![Trpdbt](https://raw.githubusercontent.com/tb0hdan/torpedo/master/doc/UDvNqA-29o.png)

Matrix: @TorpedoBot:matrix.org

IRC: #torpedobot on FreeNode



# Using Docker image

Please refer to: https://hub.docker.com/r/tb0hdan/torpedo/


# Running locally

Get Slack token(s):

`https://api.slack.com/custom-integrations/legacy-tokens`

Paste token as `token.sh`

Get Telegram/Jabber accounts.

Get Skype channel creds (https://dev.botframework.com/)

Get Sentry.io DSN: https://sentry.io

Optional parameters (all or any combination of)

```bash
SLACK="xxxttt,aaabbb"
TELEGRAM="xxx,yyy"
JABBER="user@host.com:supersecret,user2@anotherhost.com:a1FvH12"
SKYPE="app_id:app_password,app_id2:app_password2"
SENTRY_DSN="https://xxx:yyy"
FACEBOOK="aaabbb:ccc"
KIK="ddd:eee"
LINE="chat_secret:chat_token"
MATRIX="aaa:MDAxxxxxxxxxxxxxxxxxxxxx"
```


Mandatory parameters:


```bash
LASTFM_KEY="aaa"
LASTFM_SECRET="bbb"
GOOGLE_WEBAPP_KEY="ccc"
PINTEREST="ddd"
```

# Requirements

An accessible MongoDB instance (defaults to localhost)

Unauthenticated access (default):


`torpedobot -mongo host` or `torpedobot -mongo host:port`


Authenticated access:


`torpedobot -mongo mongodb://user:pass@host:port`



# Running

```bash
make deps
./run.sh
```

# Commands

## Command Prefix

Slack: `!`

Telegram: `/`

Jabber: `!`

Skype: `!` or @Botname `!`

Facebook: `!`

Kik: `!`

Line: `!`

Matrix: `!`

## Help

P stands for prefix above

```
P?
Ph
Phelp
```


e.g. for Slack it's `!help`

# Additional topics

## [TRPE](doc/TRPE.md)
## [Blacklist functionality](doc/BLACKLIST.md)
## [Development](doc/Development.md)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftb0hdan%2Ftorpedo.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftb0hdan%2Ftorpedo?ref=badge_large)