package main

import (
        "github.com/nlopes/slack"
       )

func RegisterChatHandlers(handlers map[string]func(*slack.Client, *slack.MessageEvent)) {
        handlers["bash"] = BashProcessMessage
        handlers["lastfm"] = LastFmProcessMessage
}
