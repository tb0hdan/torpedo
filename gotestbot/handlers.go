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
        // Probably useless stuff
        handlers["b64e"] = CryptoProcessMessage
        handlers["b64d"] = CryptoProcessMessage
        handlers["md5"] = CryptoProcessMessage
        handlers["sha1"] = CryptoProcessMessage
        handlers["sha256"] = CryptoProcessMessage
        handlers["sha512"] = CryptoProcessMessage
}
