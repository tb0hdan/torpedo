package multibot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	common "github.com/tb0hdan/torpedo_common"
	database "github.com/tb0hdan/torpedo_common/database"
	memcache "github.com/tb0hdan/torpedo_common/memcache"
	"github.com/tb0hdan/torpedo_registry"

	"github.com/getsentry/raven-go"
)

var bot *TorpedoBot
var once sync.Once

type BotStats struct {
	StartTimestamp         int64
	ProcessedMessages      int64
	ProcessedMessagesTotal int64
	ConnectedAccounts      int32
	TotalAccounts          int32
}

type TorpedoBot struct {
	caches          map[string]*memcache.MemCacheType
	commandHandlers map[string]func(*torpedo_registry.BotAPI, interface{}, string)
	help            map[string]string
	Database        *database.MongoDB
	logger              *log.Logger
	throttle            *memcache.MemCacheType
	RegisteredProtocols map[string]func(interface{}, string, *TorpedoBotAPI, []torpedo_registry.RichMessage)
	Stats               BotStats
	Build               struct {
		Build     string
		BuildDate string
		Version   string
	}
}

type TorpedoBotAPI struct {
	API           interface{}
	CommandPrefix string
	Bot           *TorpedoBot
	From          string
	Type          string
}

func (tba *TorpedoBotAPI) PostMessage(channel interface{}, message string, richmsgs ...torpedo_registry.RichMessage) {
	ran := 0
	for proto := range tba.Bot.RegisteredProtocols {
		if proto == fmt.Sprintf("%T", tba.API) {
			ran += 1
			tba.Bot.RegisteredProtocols[proto](channel, message, tba, richmsgs)
			break
		}
	}
	if ran == 0 {
		tba.Bot.logger.Printf("Unsupported bot API: %T\n", tba.API)

	}
}

func (tb *TorpedoBot) PostMessage(channel interface{}, message string, api *torpedo_registry.BotAPI, richmsgs ...interface{}) {
	mapi := api.API.(*TorpedoBotAPI)
	if len(richmsgs) > 0 {
		mapi.PostMessage(channel, message, richmsgs[0].(torpedo_registry.RichMessage))
	} else {
		mapi.PostMessage(channel, message)
	}
}

func (tb *TorpedoBot) processChannelEvent(api *TorpedoBotAPI, channel interface{}, incoming_message string) {
	if strings.HasPrefix(incoming_message, api.CommandPrefix) && tb.NoSpam(channel, incoming_message) {
		tb.Stats.ProcessedMessages += 1
		// is it good idea to store it here?
		// TODO: find better way
		tb.Stats.ProcessedMessagesTotal = tb.Database.GetUpdateTotalMessages(1)
		//
		command := strings.TrimPrefix(incoming_message, api.CommandPrefix)
		botapi := &torpedo_registry.BotAPI{}
		botapi.API = api
		botapi.CommandPrefix = api.CommandPrefix
		botapi.Bot.GetCachedItem = api.Bot.GetCachedItem
		botapi.Bot.SetCachedItems = api.Bot.SetCachedItems
		botapi.Bot.GetCommandHandlers = api.Bot.GetCommandHandlers
		botapi.Bot.GetHelp = api.Bot.GetHelp
		botapi.Bot.PostMessage = api.Bot.PostMessage
		botapi.Bot.Stats = api.Bot.Stats
		botapi.Bot.Build = api.Bot.Build
		found := 0
		for handler := range tb.commandHandlers {
			if strings.ToLower(strings.Split(command, " ")[0]) == handler {
				found += 1
				if torpedo_registry.Config.GetConfig()["raven"] == "yes" {
					raven.CapturePanicAndWait(func() {
						tb.commandHandlers[handler](botapi, channel, incoming_message)
					}, nil)
				} else {
					tb.commandHandlers[handler](botapi, channel, incoming_message)
				}
				break
			}
		}
		tb.logger.Printf("PROCESS! -> `%s`", command)
		if found == 0 {
			api.PostMessage(channel, fmt.Sprintf("Could not process your message: %s%s. Command unknown. Send `%shelp` for list of valid commands and `%shelp command` for details.", api.CommandPrefix, command, api.CommandPrefix, api.CommandPrefix))
		}
	}
}

func (tb *TorpedoBot) RunLoop() {
	for {
		time.Sleep(time.Second)
	}
}

func (tb *TorpedoBot) RunBotsCSV(method func(apiKey, cmd_prefix string), CSV, cmd_prefix string) {
	wrapped := func(a, b string) {}
	if torpedo_registry.Config.GetConfig()["raven"] == "yes" {
		wrapped = func(apiKey, cmd_prefix string) {
			// this should (!) capture bot protocol panic
			raven.CapturePanicAndWait(func() {
				method(apiKey, cmd_prefix)
			}, nil)
		}
	} else {
		wrapped = method
	}
	for _, key := range strings.Split(CSV, ",") {
		if key == "" {
			continue
		}
		tb.Stats.TotalAccounts += 1
		go wrapped(key, cmd_prefix)
	}
}

func (tb *TorpedoBot) RegisterHandlers(handlers map[string]func(*torpedo_registry.BotAPI, interface{}, string)) {
	tb.commandHandlers = handlers
	return
}

func (tb *TorpedoBot) GetCommandHandlers() (handlers map[string]func(*torpedo_registry.BotAPI, interface{}, string)) {
	return tb.commandHandlers
}

func (tb *TorpedoBot) RegisterHelp(help map[string]string) {
	tb.help = help
	return
}

func (tb *TorpedoBot) GetHelp() (help map[string]string) {
	return tb.help
}

func (tb *TorpedoBot) SetBuildInfo(build, buildDate, version string) {
	tb.Build.Build = build
	tb.Build.BuildDate = buildDate
	tb.Build.Version = version
	return
}

func (tb *TorpedoBot) RunPreParsers() {
	for pname, pfunc := range torpedo_registry.Config.GetPreParsers() {
		tb.logger.Printf("Running argument preparser %s\n", pname)
		pfunc(torpedo_registry.Config)
	}
	return
}

func (tb *TorpedoBot) RunPostParsers() {
	for pname, pfunc := range torpedo_registry.Config.GetPostParsers() {
		tb.logger.Printf("Running argument postparser %s\n", pname)
		pfunc(torpedo_registry.Config)
	}
	return
}

func New() *TorpedoBot {
	once.Do(func() {
		bot = &TorpedoBot{}
		cu := &common.Utils{}
		bot.logger = cu.NewLog("torpedo-bot")
		bot.caches = make(map[string]*memcache.MemCacheType)
		bot.throttle = memcache.New()
		env_dsn := os.Getenv("SENTRY_DSN")
		if env_dsn != "" {
			bot.logger.Print("Using Sentry error reporting...\n")
			raven.SetDSN(env_dsn)
			torpedo_registry.Config.SetConfig("raven", "yes")
		}
		bot.RegisteredProtocols = make(map[string]func(interface{}, string, *TorpedoBotAPI, []torpedo_registry.RichMessage))
		bot.Stats = BotStats{}
		bot.Stats.StartTimestamp = int64(time.Now().Unix())

	})
	return bot
}
