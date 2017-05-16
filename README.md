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
SLACK="xxxttt,aaabbb"
TELEGRAM="xxx,yyy"
JABBER="user@host.com:supersecret,user2@anotherhost.com:a1FvH12"
SKYPE="app_id:app_password,app_id2:app_password2"
```


Mandatory parameters:


```bash
LASTFM_KEY="aaa"
LASTFM_SECRET="bbb"
GOOGLE_WEBAPP_KEY="ccc"
PINTEREST="ddd"
```

# Running

```bash
make deps
./run.sh
```

# Commands

## Command Prefix

Slack: `!`

Telegram: `/`

Jbber: `!`

Skype: `!` or @Botname `!`


## Help

P stands for prefix above

```
P?
Ph
Phelp
```

e.g. for Slack it's `!help`
