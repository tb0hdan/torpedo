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
		help				   = make(map[string]string)
		lastfm_key             = flag.String("lastfm_key", "", "Last.FM API Key")
		lastfm_secret          = flag.String("lastfm_secret", "", "Last.FM API Secret")
		line_creds             = flag.String("line", "", "Line.Me credentials client_secret:client_token,")
		line_incoming_addr     = flag.String("line_incoming_addr", "0.0.0.0:3981", "Listen on this address for incoming Line.Me messages")
		pinterest_token        = flag.String("pinterest_token", "", "Pinterest Client Token")
		matrix                 = flag.String("matrix", "", "Matrix.org creds: ID:AccessToken,")
		mongo 		       = flag.String("mongo", "", "MongoDB server hostname")
	)
	flag.Parse()
	handlers["bashim"] = BashProcessMessage
	help["bashim"] = "Get random quote from Bash.im"
	handlers["bashorg"] = BashOrgProcessMessage
	help["bashorg"] = "Get random quote from bash.org"
	handlers["lastfm"] = LastFmProcessMessage
	help["lastfm"] = "Query Last.FM"
	handlers["qr"] = QREncoderProcessMessage
	help["qr"] = "Create QR Code from URL"
	handlers["tinyurl"] = TinyURLProcessMessage
	help["tinyurl"] = "Shorten URL using TinyURL.com"
	// in progress
	// handlers["pinterest"] = PinterestProcessMessage
	// Help handlers
	handlers["?"] = HelpProcessMessage
	help["?"] = "Get help using this command"
	handlers["h"] = HelpProcessMessage
	help["?"] = "Get help using this command"
	handlers["help"] = HelpProcessMessage
	help["?"] = "Get help using this command"
	// Multimedia
	handlers["wiki"] = WikiProcessMessage
	help["wiki"] = "Get article excerpt from Wikipedia.org"
	handlers["youtube"] = YoutubeProcessMessage
	help["youtube"] = "Get Youtube.com URL for specified track"
	// in progress
	handlers["giphy"] = GiphyProcessMessage
	help["giphy"] = "Get Giphy.com image"
	// Probably useless stuff
	handlers["b64e"] = CryptoProcessMessage
	help["b64e"] = "Base64 encode"
	handlers["b64d"] = CryptoProcessMessage
	help["b64d"] = "Base64 decode"
	handlers["md5"] = CryptoProcessMessage
	help["md5"] = "Calculate message MD5 sum"
	handlers["sha1"] = CryptoProcessMessage
	help["sha1"] = "Calculate message SHA1 sum"
	handlers["sha256"] = CryptoProcessMessage
	help["sha256"] = "Calculate message SHA256 sum"
	handlers["sha512"] = CryptoProcessMessage
	help["sha512"] = "Calculate message SHA512 sum"
	// entertainment :)
	handlers["sudo"] = FunProcessMessage
	help["sudo"] = "Run sudo on this machine"
	handlers["rm"] = FunProcessMessage
	help["rm"] = "Remove files on this machine"
	handlers["shutdown"] = FunProcessMessage
	help["shutdown"] = "Shutdown this machine for good"
	handlers["halt"] = FunProcessMessage
	help["halt"] = "Halt this machine for good"
	handlers["reboot"] = FunProcessMessage
	help["reboot"] = "Reboot this machine"
	handlers["poweroff"] = FunProcessMessage
	help["poweroff"] = "Power off this machine"
	handlers["kill"] = FunProcessMessage
	help["kill"] = "Terminate any process running on this machine"
	//
	handlers["so"] = StackOverflowProcessMessage
	help["so"] = "Search for solution on StackOverflow.com"
	handlers["stackoverflow"] = StackOverflowProcessMessage
	help["stackoverflow"] = "Search for solution on StackOverflow.com"
	// maps
	handlers["distance"] = DistanceProcessMessage
	help["distance"] = "Get driving distance between cities"
	// statistics
	handlers["stats"] = StatsProcessMessage
	help["stats"] = "Just system stats, nothing interesting"
	// File category (broken!)
	//handlers["getimg"] = GetSetImageProcessMessage
	//handlers["lsimg"] = GetSetImageProcessMessage
	//handlers["listimg"] = GetSetImageProcessMessage
	//handlers["setimg"] = GetSetImageProcessMessage
	//handlers["rmimg"] = GetSetImageProcessMessage
	// Games'n'stuff
	// needs better formatting + Facebook rework
	handlers["steam"] = SteamProcessMessage
	help["steam"] = "Get http://store.steampowered.com/ deals"

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
	if *mongo == "" {
		// try supplied one first
		*mongo = GetStripEnv("MONGO")
		// docker...
		if *mongo == "" {
			*mongo = GetStripEnv("MONGO_PORT_27017_TCP_ADDR")
		}

	}

	bot := multibot.New(*facebook_incoming_addr, *google_webapp_key,
		*skype_incoming_addr, *kik_incoming_addr,
		*kik_webhook_url,
		*lastfm_key, *lastfm_secret, *line_incoming_addr, *pinterest_token,
		*mongo)
	bot.RegisterHandlers(handlers)
	bot.RegisterHelp(help)
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
