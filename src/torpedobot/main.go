package main

import (
	"flag"
	"os"
	"strings"

	"torpedobot/multibot"
)

func GetStripEnv(envvar string) (result string) {
	result = os.Getenv(envvar)
	result = strings.TrimLeft(result, "'")
	result = strings.TrimRight(result, "'")
	return
}

func main() {
	var (
		slack                  = flag.String("slack", "", "Comma separated list of Slack legacy tokens")
		telegram               = flag.String("telegram", "", "Comma separated list of Telegram bot keys")
		jabber                 = flag.String("jabber", "", "Comma separated list of jabber creds, user@host.com:password,")
		skype                  = flag.String("skype", "", "Comma separated list of dev.botframework.com creds, app_id:app_password,")
		kik                    = flag.String("kik", "", "Comma separated list of Kik creds, username:api_key,")
		skype_incoming_addr    = flag.String("skype_incoming_addr", "0.0.0.0:3978", "Listen on this address for incoming Skype messages")
		facebook               = flag.String("facebook", "", "Comma separated list of Facebook creds, page_token1:verify_token1,..")
		google_webapp_key      = flag.String("google_webapp_key", "", "Google Data API Web Application Key")
		facebook_incoming_addr = flag.String("facebook_incoming_addr", "0.0.0.0:3979", "Listen on this address for incoming Facebook messages")
		kik_incoming_addr      = flag.String("kik_incoming_addr", "0.0.0.0:3980", "Listen on this address for incoming Kik messages")
		kik_webhook_url        = flag.String("kik_webhook_url", "", "Webhook URL (external) for incoming Kik messages")
		handlers               = make(map[string]func(*multibot.TorpedoBotAPI, interface{}, string))
		lastfm_key             = flag.String("lastfm_key", "", "Last.FM API Key")
		lastfm_secret          = flag.String("lastfm_secret", "", "Last.FM API Secret")
		line_creds             = flag.String("line", "", "Line.Me credentials client_secret:client_token,")
		line_incoming_addr     = flag.String("line_incoming_addr", "0.0.0.0:3981", "Listen on this address for incoming Line.Me messages")
		pinterest_token        = flag.String("pinterest_token", "", "Pinterest Client Token")
		matrix                 = flag.String("matrix", "", "Matrix.org creds: ID:AccessToken,")
	)
	flag.Parse()
	handlers["bashim"] = BashProcessMessage
	handlers["bashorg"] = BashOrgProcessMessage
	handlers["lastfm"] = LastFmProcessMessage
	handlers["qr"] = QREncoderProcessMessage
	handlers["tinyurl"] = TinyURLProcessMessage
	// in progress
	// handlers["pinterest"] = PinterestProcessMessage
	// Help handlers
	handlers["?"] = HelpProcessMessage
	handlers["h"] = HelpProcessMessage
	handlers["help"] = HelpProcessMessage
	// Multimedia
	handlers["wiki"] = WikiProcessMessage
	handlers["youtube"] = YoutubeProcessMessage
	// in progress
	handlers["giphy"] = GiphyProcessMessage
	// Probably useless stuff
	handlers["b64e"] = CryptoProcessMessage
	handlers["b64d"] = CryptoProcessMessage
	handlers["md5"] = CryptoProcessMessage
	handlers["sha1"] = CryptoProcessMessage
	handlers["sha256"] = CryptoProcessMessage
	handlers["sha512"] = CryptoProcessMessage
	// statistics
	handlers["stats"] = StatsProcessMessage
	// File category (broken!)
	//handlers["getimg"] = GetSetImageProcessMessage
	//handlers["lsimg"] = GetSetImageProcessMessage
	//handlers["listimg"] = GetSetImageProcessMessage
	//handlers["setimg"] = GetSetImageProcessMessage
	//handlers["rmimg"] = GetSetImageProcessMessage
	// Games'n'stuff
	// needs better formatting + Facebook rework
	handlers["steam"] = SteamProcessMessage

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
	if *kik == "" {
		*kik = GetStripEnv("KIK")
	}
	if *kik_webhook_url == "" {
		*kik_webhook_url = GetStripEnv("KIK_WEBHOOK_URL")
	}
	if *lastfm_key == "" {
		*lastfm_key = GetStripEnv("LASTFM_KEY")
	}
	if *lastfm_secret == "" {
		*lastfm_secret = GetStripEnv("LASTFM_SECRET")
	}
	if *line_creds == "" {
		*line_creds = GetStripEnv("LINE")
	}
	if *pinterest_token == "" {
		*pinterest_token = GetStripEnv("PINTEREST")
	}
	if *google_webapp_key == "" {
		*google_webapp_key = GetStripEnv("GOOGLE_WEBAPP_KEY")
	}
	if *matrix == "" {
		*matrix = GetStripEnv("MATRIX")
	}

	bot := multibot.New(*facebook_incoming_addr, *google_webapp_key,
		*skype_incoming_addr, *kik_incoming_addr,
		*kik_webhook_url,
		*lastfm_key, *lastfm_secret, *line_incoming_addr, *pinterest_token)
	bot.RegisterHandlers(handlers)
	bot.RunBotsCSV(bot.RunSlackBot, *slack, "!")
	bot.RunBotsCSV(bot.RunTelegramBot, *telegram, "/")
	bot.RunBotsCSV(bot.RunJabberBot, *jabber, "!")
	bot.RunBotsCSV(bot.RunSkypeBot, *skype, "!")
	bot.RunBotsCSV(bot.RunFacebookBot, *facebook, "!")
	bot.RunBotsCSV(bot.RunKikBot, *kik, "!")
	bot.RunBotsCSV(bot.RunLineBot, *line_creds, "!")
	bot.RunBotsCSV(bot.RunMatrixBot, *matrix, "!")
	bot.RunLoop()
}
