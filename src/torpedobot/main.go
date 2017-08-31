package main

import (
	"flag"
	"fmt"
	"os"

	"torpedobot/multibot"

	"github.com/erikdubbelboer/gspt"
	"github.com/tb0hdan/torpedo_registry"
)

const ProjectURL = "https://github.com/tb0hdan/torpedo"

// Global vars for versioning
var (
	BUILD      = "Not available"
	BUILD_DATE = "Not available"
	VERSION    = "Not available"
)

func main() {
	// Set process name - https://github.com/tb0hdan/torpedo/issues/30
	gspt.SetProcTitle(os.Args[0])

	// Help handlers
	help_msg := "Get help using this command"
	torpedo_registry.Config.RegisterHelpAndHandler("?", help_msg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("h", help_msg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("help", help_msg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("stats", "Just system stats, nothing interesting", StatsProcessMessage)

	bot := multibot.New()
	bot.SetBuildInfo(BUILD, BUILD_DATE, VERSION, ProjectURL)
	// bot cfg
	// plugins/protocols
	torpedo_registry.Config.RegisterParser("slack", bot.ConfigureSlackBot, bot.ParseSlackBot)
	torpedo_registry.Config.RegisterParser("telegram", bot.ConfigureTelegramBot, bot.ParseTelegramBot)
	torpedo_registry.Config.RegisterParser("jabber", bot.ConfigureJabberBot, bot.ParseJabberBot)
	torpedo_registry.Config.RegisterParser("skype", bot.ConfigureSkypeBot, bot.ParseSkypeBot)
	torpedo_registry.Config.RegisterParser("kik", bot.ConfigureKikBot, bot.ParseKikBot)
	torpedo_registry.Config.RegisterParser("line", bot.ConfigureLineBot, bot.ParseLineBot)
	torpedo_registry.Config.RegisterParser("matrix", bot.ConfigureMatrixBot, bot.ParseMatrixBot)
	torpedo_registry.Config.RegisterParser("facebook", bot.ConfigureFacebookBot, bot.ParseFacebookBot)
	torpedo_registry.Config.RegisterParser("irc", bot.ConfigureIRCBot, bot.ParseIRCBot)

	// internals
	torpedo_registry.Config.RegisterParser("apiaddr", bot.ConfigureHTTPAPI, bot.ParseHTTPAPI)
	torpedo_registry.Config.RegisterParser("mongodb", bot.ConfigureMongoDBPlugin, bot.ParseMongoDBPlugin)
	torpedo_registry.Config.RegisterParser("trpe", bot.ConfigureTRPE, bot.ParseTRPE)

	bot.RunPreParsers()
	flag.Parse()
	bot.RunPostParsers()

	fmt.Println(torpedo_registry.Config.GetConfig())
	bot.RunBotsCSV(bot.RunSlackBot, torpedo_registry.Config.GetConfig()["slackapikey"], "!")
	bot.RunBotsCSV(bot.RunTelegramBot, torpedo_registry.Config.GetConfig()["telegramapikey"], "/")
	bot.RunBotsCSV(bot.RunJabberBot, torpedo_registry.Config.GetConfig()["jabberapikey"], "!")
	bot.RunBotsCSV(bot.RunSkypeBot, torpedo_registry.Config.GetConfig()["skypeapikey"], "!")
	bot.RunBotsCSV(bot.RunKikBot, torpedo_registry.Config.GetConfig()["kikapikey"], "!")
	bot.RunBotsCSV(bot.RunLineBot, torpedo_registry.Config.GetConfig()["lineapikey"], "!")
	bot.RunBotsCSV(bot.RunMatrixBot, torpedo_registry.Config.GetConfig()["matrixapikey"], "!")
	bot.RunBotsCSV(bot.RunFacebookBot, torpedo_registry.Config.GetConfig()["facebookapikey"], "!")
	bot.RunBotsCSV(bot.RunIRCBot, torpedo_registry.Config.GetConfig()["ircapikey"], "!")

	// start plugin coroutines (if any) after connecting to accounts
	bot.RunCoroutines()
	// start HTTP API Server
	bot.RunHTTPAPI()
	// start eternal loop
	bot.RunLoop()
}
