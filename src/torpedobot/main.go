package main

import (
	"flag"
	"strings"

	"torpedobot/multibot"

)

var (
	token    = flag.String("token", "", "Comma separated list of Slack legacy tokens")
	telegram = flag.String("telegram", "", "Comma separated list of Telegram bot keys")
	jabber = flag.String("jabber", "", "Comma separated list of jabber creds, user@host.com:password,")
	handlers = make(map[string]func(*multibot.TorpedoBotAPI, *multibot.TorpedoBot, interface{}, string, string))
)

func main() {
	var keys []string
	var telegram_keys []string
	var jabber_keys []string

	flag.Parse()

	for _, key := range strings.Split(*token, ",") {
		keys = append(keys, key)
	}

	for _, key := range strings.Split(*telegram, ",") {
		telegram_keys = append(telegram_keys, key)
	}

	for _, key := range strings.Split(*jabber, ",") {
		jabber_keys = append(jabber_keys, key)
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

	bot := multibot.New()
	bot.RegisterHandlers(handlers)
	bot.RunSlackBots(keys, "!")
	bot.RunTelegramBots(telegram_keys, "/")
	bot.RunJabberBots(jabber_keys, "!")
	bot.RunLoop()
}
