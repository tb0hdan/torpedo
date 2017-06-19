package main

import (
	"flag"

	"torpedobot/multibot"
)

// Global vars for versioning
var (
	BUILD      = "Not available"
	BUILD_DATE = "Not available"
	VERSION    = "Not available"
)

func main() {

	var (
		handlers = make(map[string]func(*multibot.TorpedoBotAPI, interface{}, string))
		help     = make(map[string]string)
	)
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
	handlers["soundcloud"] = SoundCloudProcessMessage
	handlers["sc"] = SoundCloudProcessMessage
	help["soundcloud"] = "Search for track on SoundCloud"
	help["sc"] = "Search for track on SoundCloud"
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
	//
	handlers["xkcd"] = XKCDProcessMessage
	help["xkcd"] = "Get XKCD random strip. Provide integer ID to get specific one."

	bot := multibot.New()
	bot.SetBuildInfo(BUILD, BUILD_DATE, VERSION)
	//
	bot.ConfigureSlackBot()
	bot.ConfigureTelegramBot()
	bot.ConfigureJabberBot()
	bot.ConfigureSkypeBot()
	bot.ConfigureKikBot()
	bot.ConfigureFacebookBot()
	bot.ConfigureLineBot()
	// bot plugins
	bot.ConfigureLastFmPlugin()
	bot.ConfigureGooglePlugin()
	bot.ConfigurePinterestPlugin()
	bot.ConfigureSoundCloudPlugin()
	bot.ConfigureMongoDBPlugin()
	// make this one last
	flag.Parse()

	// bot plugins
	bot.RunLastFmPlugin()
	bot.RunGooglePlugin()
	bot.RunPinterestPlugin()
	bot.RunSoundCloudPlugin()
	bot.RunMongoDBPlugin()
	//
	bot.RegisterHandlers(handlers)
	bot.RegisterHelp(help)
	bot.RunBotsCSV(bot.RunSlackBot, bot.Config.SlackAPIKey, "!")
	bot.RunBotsCSV(bot.RunTelegramBot, bot.Config.TelegramAPIKey, "/")
	bot.RunBotsCSV(bot.RunJabberBot, bot.Config.JabberAPIKey, "!")
	bot.RunBotsCSV(bot.RunSkypeBot, bot.Config.SkypeAPIKey, "!")
	bot.RunBotsCSV(bot.RunFacebookBot, bot.Config.FacebookAPIKey, "!")
	bot.RunBotsCSV(bot.RunKikBot, bot.Config.KikAPIKey, "!")
	bot.RunBotsCSV(bot.RunLineBot, bot.Config.LineAPIKey, "!")
	bot.RunBotsCSV(bot.RunMatrixBot, bot.Config.MatrixAPIKey, "!")
	bot.RunLoop()
}
