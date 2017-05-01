package main

import (
	"github.com/nlopes/slack"
)

func RegisterChatHandlers()  (handlers map[string]func(*slack.Client, *slack.MessageEvent)) {
	handlers = make(map[string]func(*slack.Client, *slack.MessageEvent))
	// TODO: Support command groups and provide category help
	handlers["bashim"] = BashProcessMessage
	handlers["bashorg"] = BashOrgProcessMessage
 	handlers["lastfm"] = LastFmProcessMessage
	handlers["qr"] = QREncoderProcessMessage
	handlers["tinyurl"] = TinyURLProcessMessage
	handlers["pinterest"] = PinterestProcessMessage
	// Help handlers
	handlers["?"] = HelpProcessMessage
	handlers["h"] = HelpProcessMessage
	handlers["help"] = HelpProcessMessage
	// Multimedia
	handlers["wiki"] = WikiProcessMessage
	handlers["youtube"] = YoutubeProcessMessage
	handlers["giphy"] = GiphyProcessMessage
	// Probably useless stuff
	handlers["b64e"] = CryptoProcessMessage
	handlers["b64d"] = CryptoProcessMessage
	handlers["md5"] = CryptoProcessMessage
	handlers["sha1"] = CryptoProcessMessage
	handlers["sha256"] = CryptoProcessMessage
	handlers["sha512"] = CryptoProcessMessage
	// File category
	handlers["getimg"] = GetSetImageProcessMessage
	handlers["lsimg"] = GetSetImageProcessMessage
	handlers["listimg"] = GetSetImageProcessMessage
	handlers["setimg"] = GetSetImageProcessMessage
	handlers["rmimg"] = GetSetImageProcessMessage
	return
}
