package main

import (
	"flag"
	"os"
	"strings"

	"torpedobot/multibot"
)


var (
	slack               = flag.String("slack", "", "Comma separated list of Slack legacy tokens")
	telegram            = flag.String("telegram", "", "Comma separated list of Telegram bot keys")
	jabber              = flag.String("jabber", "", "Comma separated list of jabber creds, user@host.com:password,")
	skype               = flag.String("skype", "", "Comma separated list of dev.botframework.com creds, app_id:app_password,")
	skype_incoming_addr = flag.String("skype_incoming_addr", "0.0.0.0:3978", "Listen on this address for incoming Skype messages")
	facebook = flag.String("facebook", "", "Comma separated list of Facebook creds, page_token1:verify_token1,..")
	facebook_incoming_addr = flag.String("facebook_incoming_addr", "0.0.0.0:3979", "Listen on this address for incoming Facebook messages")
	handlers            = make(map[string]func(*multibot.TorpedoBotAPI, interface{}, string))
)


func GetStripEnv(envvar string) (result string) {
	result = os.Getenv(envvar)
	result = strings.TrimLeft(result, "'")
	result = strings.TrimRight(result, "'")
	return
}


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

	bot := multibot.New(*facebook_incoming_addr, *skype_incoming_addr)
	bot.RegisterHandlers(handlers)
	if *slack == "" {
		*slack = GetStripEnv("SLACK")
	}
	if *telegram == "" {
		*telegram = GetStripEnv("TELEGRAM")
	}
	if *jabber == "" {
		*jabber = GetStripEnv("JABBER")
	}
	if *skype == "" {
		*skype = GetStripEnv("SKYPE")
	}
	if *facebook == "" {
		*facebook = GetStripEnv("FACEBOOK")
	}
	if *lastfm_key == "" {
		*lastfm_key = GetStripEnv("LASTFM_KEY")
	}
	if *lastfm_secret == "" {
		*lastfm_secret = GetStripEnv("LASTFM_SECRET")
	}
	if *pinterest_token == "" {
		*pinterest_token = GetStripEnv("PINTEREST")
	}
	if *google_webapp_key == "" {
		*google_webapp_key = GetStripEnv("GOOGLE_WEBAPP_KEY")
	}
	bot.RunBotsCSV(bot.RunSlackBot, *slack, "!")
	bot.RunBotsCSV(bot.RunTelegramBot, *telegram, "/")
	bot.RunBotsCSV(bot.RunJabberBot, *jabber, "!")
	bot.RunBotsCSV(bot.RunSkypeBot, *skype, "!")
	bot.RunBotsCSV(bot.RunFacebookBot, *facebook, "!")
	bot.RunLoop()
}
