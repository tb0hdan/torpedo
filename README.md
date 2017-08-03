# torpedo
Multiprotocol bot

# Intro

Torpedo uses multiple accounts (at least one is required). Supported transports:

- Facebook
- Jabber
- Kik
- Line
- Matrix (matrix.org only atm)
- Skype
- Slack
- Telegram


# See it in action

Facebook: https://www.facebook.com/TorpedoBot/

Jabber: torpedobot@jabber.cz

Skype: https://join.skype.com/bot/f61c6815-438d-4795-8aaa-9b1d8d2a342a

Telegram: http://t.me/TorpedoTelegramBot

Line:

![Trpdbt](https://raw.githubusercontent.com/tb0hdan/torpedo/master/doc/UDvNqA-29o.png)

Matrix: @TorpedoBot:matrix.org


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



# TorpedoBot Remote Plugin Execution (TRPE)

TRPE allows writing plugins in any language as long as content is returned via HTTP API.
Sample application is available at `tools/trpe_server.py`

Architecture is as follows:

![TRPE](https://raw.githubusercontent.com/tb0hdan/torpedo/master/doc/TRPE.png)
