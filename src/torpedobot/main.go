package main


import (
	"flag"

	"torpedobot/multibot"

)


var (
	token    = flag.String("token", "", "Comma separated list of Slack legacy tokens")
	telegram = flag.String("telegram", "", "Comma separated list of Telegram bot keys")
	jabber = flag.String("jabber", "", "Comma separated list of jabber creds, user@host.com:password,")
	skype = flag.String("skype", "", "Comma separated list of dev.botframework.com creds, app_id:app_password,")
	skype_incoming_addr = flag.String("skype_incoming_addr", "localhost:3978", "Listen on this address for incoming Skype messages")
	handlers = make(map[string]func(*multibot.TorpedoBotAPI, interface{}, string))
)


func main() {
	flag.Parse()
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
	// Games'n'stuff
	handlers["steam"] = SteamProcessMessage

	bot := multibot.New(*skype_incoming_addr)
	bot.RegisterHandlers(handlers)
	bot.RunBotsCSV(bot.RunSlackBot, *token, "!")
	bot.RunBotsCSV(bot.RunTelegramBot, *telegram, "/")
	bot.RunBotsCSV(bot.RunJabberBot, *jabber, "!")
	bot.RunBotsCSV(bot.RunSkypeBot, *skype, "!")
	bot.RunLoop()
}
