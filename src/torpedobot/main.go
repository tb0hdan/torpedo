package main

import (
	"flag"
	"fmt"

	"torpedobot/multibot"

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
	// Help handlers
	help_msg := "Get help using this command"
	torpedo_registry.Config.RegisterHelpAndHandler("?", help_msg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("h", help_msg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("help", help_msg, HelpProcessMessage)
	torpedo_registry.Config.RegisterHelpAndHandler("stats", "Just system stats, nothing interesting", StatsProcessMessage)

	bot := multibot.New()
	bot.SetBuildInfo(BUILD, BUILD_DATE, VERSION, ProjectURL)
	// bot cfg
	torpedo_registry.Config.RegisterParser("slack", bot.ConfigureSlackBot, bot.ParseSlackBot)
	torpedo_registry.Config.RegisterParser("telegram", bot.ConfigureTelegramBot, bot.ParseTelegramBot)
	torpedo_registry.Config.RegisterParser("jabber", bot.ConfigureJabberBot, bot.ParseJabberBot)
	torpedo_registry.Config.RegisterParser("skype", bot.ConfigureSkypeBot, bot.ParseSkypeBot)
	torpedo_registry.Config.RegisterParser("kik", bot.ConfigureKikBot, bot.ParseKikBot)
	torpedo_registry.Config.RegisterParser("line", bot.ConfigureLineBot, bot.ParseLineBot)
	torpedo_registry.Config.RegisterParser("matrix", bot.ConfigureMatrixBot, bot.ParseMatrixBot)
	torpedo_registry.Config.RegisterParser("facebook", bot.ConfigureFacebookBot, bot.ParseFacebookBot)
	torpedo_registry.Config.RegisterParser("mongodb", bot.ConfigureMongoDBPlugin, bot.ParseMongoDBPlugin)

	bot.RunPreParsers()
	flag.Parse()
	bot.RunPostParsers()

	// Command handlers and help
	bot.RegisterHandlers(torpedo_registry.Config.GetHandlers())
	bot.RegisterHelp(torpedo_registry.Config.GetHelp())

	fmt.Println(torpedo_registry.Config.GetConfig())
	bot.RunBotsCSV(bot.RunSlackBot, torpedo_registry.Config.GetConfig()["slackapikey"], "!")
	bot.RunBotsCSV(bot.RunTelegramBot, torpedo_registry.Config.GetConfig()["telegramapikey"], "/")
	bot.RunBotsCSV(bot.RunJabberBot, torpedo_registry.Config.GetConfig()["jabberapikey"], "!")
	bot.RunBotsCSV(bot.RunSkypeBot, torpedo_registry.Config.GetConfig()["skypeapikey"], "!")
	bot.RunBotsCSV(bot.RunKikBot, torpedo_registry.Config.GetConfig()["kikapikey"], "!")
	bot.RunBotsCSV(bot.RunLineBot, torpedo_registry.Config.GetConfig()["lineapikey"], "!")
	bot.RunBotsCSV(bot.RunMatrixBot, torpedo_registry.Config.GetConfig()["matrixapikey"], "!")
	bot.RunBotsCSV(bot.RunFacebookBot, torpedo_registry.Config.GetConfig()["facebookapikey"], "!")

	// start plugin coroutines (if any) after connecting to accounts
	bot.RunCoroutines()
	bot.RunLoop()
}
