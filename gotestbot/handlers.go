package main

import (
        "github.com/nlopes/slack"
       )

func RegisterChatHandlers(handlers map[string]func(*slack.Client, *slack.MessageEvent)) {
        handlers["bash"] = BashProcessMessage
        handlers["lastfm"] = LastFmProcessMessage
        handlers["qr"] = QREncoderProcessMessage
        // Help handlers
        handlers["?"] = HelpProcessMessage
        handlers["h"] = HelpProcessMessage
        handlers["help"] = HelpProcessMessage
        //
        handlers["wiki"] = WikiProcessMessage
        handlers["youtube"] = YoutubeProcessMessage
}
