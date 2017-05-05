package main

import (
	"flag"
	"strings"

	"github.com/nlopes/slack"

	"torpedobot/multibot"
)

var (
	token = flag.String("token", "", "Comma separated list of Slack legacy tokens")
	handlers = make(map[string]func(*slack.Client, *slack.MessageEvent, *multibot.TorpedoBot))
)

func main() {
	var keys []string

	flag.Parse()

	for _, key := range strings.Split(*token, ",") {
		keys = append(keys, key)
	}

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

	bot := multibot.New(keys)
	bot.RegisterHandlers(handlers)
	bot.RunSlackBots()
	bot.RunLoop()
}
