package main

import (
	"flag"
	"os"

	"torpedobot/multibot"

	"github.com/erikdubbelboer/gspt"
	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
)

// Global vars for versioning
var (
	Build      = "Not available"
	BuildDate  = "Not available"
	GoVersion  = "Not available"
	Version    = "Not available"
	ProjectURL = "Not available"
)

func main() {
	// Set process name - https://github.com/tb0hdan/torpedo/issues/30
	gspt.SetProcTitle(os.Args[0])

	// Help handlers
	helpMsg := "Get help using this command"
	torpedo_registry.Config.RegisterHelpAndHandler("?", helpMsg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("h", helpMsg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("help", helpMsg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("stats", "Just system stats, nothing interesting", StatsProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("chatinfo", "Chat/DM information", ChatInfoProcessMessage)

	bot := multibot.New()
	bot.SetBuildInfo(Build, BuildDate, GoVersion, Version, ProjectURL)
	// bot cfg
	// plugins/protocols
	torpedo_registry.Config.RegisterParser("slack", bot.ConfigureSlackBot, bot.ParseSlackBot)
	torpedo_registry.Config.RegisterParser("telegram", bot.ConfigureTelegramBot, bot.ParseTelegramBot)
	torpedo_registry.Config.RegisterParser("jabber", bot.ConfigureJabberBot, bot.ParseJabberBot)
	torpedo_registry.Config.RegisterParser("skype", bot.ConfigureSkypeBot, bot.ParseSkypeBot)
	torpedo_registry.Config.RegisterParser("teams", bot.ConfigureTeamsBot, bot.ParseTeamsBot)
	torpedo_registry.Config.RegisterParser("kik", bot.ConfigureKikBot, bot.ParseKikBot)
	torpedo_registry.Config.RegisterParser("line", bot.ConfigureLineBot, bot.ParseLineBot)
	torpedo_registry.Config.RegisterParser("matrix", bot.ConfigureMatrixBot, bot.ParseMatrixBot)
	torpedo_registry.Config.RegisterParser("facebook", bot.ConfigureFacebookBot, bot.ParseFacebookBot)
	torpedo_registry.Config.RegisterParser("irc", bot.ConfigureIRCBot, bot.ParseIRCBot)

	// internals
	torpedo_registry.Config.RegisterParser("debug", bot.ConfigureDebug, bot.ParseDebug)
	torpedo_registry.Config.RegisterParser("apiaddr", bot.ConfigureHTTPAPI, bot.ParseHTTPAPI)
	torpedo_registry.Config.RegisterParser("mongodb", bot.ConfigureMongoDBPlugin, bot.ParseMongoDBPlugin)
	torpedo_registry.Config.RegisterParser("trpe", bot.ConfigureTRPE, bot.ParseTRPE)
	torpedo_registry.Config.RegisterParser("list_handlers", bot.ConfigureListPlugins, bot.ParseListPlugins)

	bot.RunPreParsers()
	flag.Parse()
	bot.RunPostParsers()

	if torpedo_registry.Config.GetConfig()["debug"] == "yes" {
		cu := &common.Utils{}
		logger := cu.NewLog("torpedo-bot")
		logger.Println(torpedo_registry.Config.GetConfig())
	}
	bot.RunBotsCSV(bot.RunSlackBot, torpedo_registry.Config.GetConfig()["slackapikey"], "!")
	bot.RunBotsCSV(bot.RunTelegramBot, torpedo_registry.Config.GetConfig()["telegramapikey"], "/")
	bot.RunBotsCSV(bot.RunJabberBot, torpedo_registry.Config.GetConfig()["jabberapikey"], "!")
	bot.RunBotsCSV(bot.RunSkypeBot, torpedo_registry.Config.GetConfig()["skypeapikey"], "!")
	bot.RunBotsCSV(bot.RunTeamsBot, torpedo_registry.Config.GetConfig()["teamsapikey"], "!")
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
