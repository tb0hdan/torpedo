package main

import (
	"flag"

	_ "github.com/tb0hdan/torpedo_plugins/torpedo_bashim_plugin"
	_ "github.com/tb0hdan/torpedo_plugins/torpedo_bashorg_plugin"
	_ "github.com/tb0hdan/torpedo_plugins/torpedo_fun_plugin"
	_ "github.com/tb0hdan/torpedo_plugins/torpedo_giphy_plugin"

	"github.com/tb0hdan/torpedo_registry"
	"torpedobot/multibot"
)

// Global vars for versioning
var (
	BUILD      = "Not available"
	BUILD_DATE = "Not available"
	VERSION    = "Not available"
)

func main() {
	// Help handlers
	help_msg := "Get help using this command"
	torpedo_registry.RegisterHandler("?", HelpProcessMessage)
	torpedo_registry.RegisterHelp("?", help_msg)
	torpedo_registry.RegisterHandler("h", HelpProcessMessage)
	torpedo_registry.RegisterHelp("h", help_msg)
	torpedo_registry.RegisterHandler("help", HelpProcessMessage)
	torpedo_registry.RegisterHelp("help", help_msg)
	torpedo_registry.RegisterHandler("stats", StatsProcessMessage)
	torpedo_registry.RegisterHelp("stats", "Just system stats, nothing interesting")

	bot := multibot.New()
	bot.SetBuildInfo(BUILD, BUILD_DATE, VERSION)
	// bot cfg
	torpedo_registry.RegisterPreParser("slack", bot.ConfigureSlackBot)
	torpedo_registry.RegisterPreParser("telegram", bot.ConfigureTelegramBot)
	torpedo_registry.RegisterPreParser("jabber", bot.ConfigureJabberBot)
	torpedo_registry.RegisterPreParser("skype", bot.ConfigureSkypeBot)
	torpedo_registry.RegisterPreParser("kik", bot.ConfigureKikBot)
	torpedo_registry.RegisterPreParser("facebook", bot.ConfigureFacebookBot)
	torpedo_registry.RegisterPreParser("line", bot.ConfigureLineBot)
	torpedo_registry.RegisterPreParser("matrix", bot.ConfigureMatrixBot)
	torpedo_registry.RegisterPreParser("mongodb", bot.ConfigureMongoDBPlugin)

	bot.RunPreParsers(torpedo_registry.GetPreParsers())

	flag.Parse()

	torpedo_registry.RegisterPostParser("facebook", bot.ParseFacebookBot)
	torpedo_registry.RegisterPostParser("jabber", bot.ParseJabberBot)
	torpedo_registry.RegisterPostParser("kik", bot.ParseKikBot)
	torpedo_registry.RegisterPostParser("line", bot.ParseLineBot)
	torpedo_registry.RegisterPostParser("matrix", bot.ParseMatrixBot)
	torpedo_registry.RegisterPostParser("skype", bot.ParseSkypeBot)
	torpedo_registry.RegisterPostParser("slack", bot.ParseSlackBot)
	torpedo_registry.RegisterPostParser("telegram", bot.ParseTelegramBot)
	torpedo_registry.RegisterPostParser("mongodb", bot.ParseMongoDBPlugin)

	bot.RunPostParsers(torpedo_registry.GetPostParsers())

	// Command handlers and help
	bot.RegisterHandlers(torpedo_registry.GetHandlers())
	bot.RegisterHelp(torpedo_registry.GetHelp())

	bot.RunBotsCSV(bot.RunSlackBot, *bot.Config.SlackAPIKey, "!")
	bot.RunBotsCSV(bot.RunTelegramBot, *bot.Config.TelegramAPIKey, "/")
	bot.RunBotsCSV(bot.RunJabberBot, *bot.Config.JabberAPIKey, "!")
	bot.RunBotsCSV(bot.RunSkypeBot, *bot.Config.SkypeAPIKey, "!")
	bot.RunBotsCSV(bot.RunFacebookBot, *bot.Config.FacebookAPIKey, "!")
	bot.RunBotsCSV(bot.RunKikBot, *bot.Config.KikAPIKey, "!")
	bot.RunBotsCSV(bot.RunLineBot, *bot.Config.LineAPIKey, "!")
	bot.RunBotsCSV(bot.RunMatrixBot, *bot.Config.MatrixAPIKey, "!")
	bot.RunLoop()
}
