package multibot

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"os/signal"
	"syscall"

	"github.com/getsentry/raven-go"
	common "github.com/tb0hdan/torpedo_common"
	database "github.com/tb0hdan/torpedo_common/database"
	memcache "github.com/tb0hdan/torpedo_common/memcache"
	"github.com/tb0hdan/torpedo_registry"
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
	caches              map[string]*memcache.MemCacheType
	Database            *database.MongoDB
	logger              *log.Logger
	throttle            *memcache.MemCacheType
	RegisteredProtocols map[string]func(interface{}, string, *TorpedoBotAPI, []torpedo_registry.RichMessage)
	Stats               BotStats
	Build               struct {
		Build      string
		BuildDate  string
		GoVersion  string
		Version    string
		ProjectURL string
	}
}

type TorpedoBotAPI struct {
	API           interface{}
	CommandPrefix string
	Bot           *TorpedoBot
	// FIXME: Move From field to UserProfile struct
	From        string
	Type        string
	UserProfile *torpedo_registry.UserProfile
	Me          string
}

// This is required for plugins to have loose coupling with bot itself
// they use torpedo_registry.BotAPI instead
func (tb *TorpedoBot) GetBotAPI(api *TorpedoBotAPI, channel interface{}, incoming_message string) (botapi *torpedo_registry.BotAPI) {
	botapi = &torpedo_registry.BotAPI{}
	botapi.API = api
	botapi.ProtocolAPI = api.API
	botapi.CommandPrefix = api.CommandPrefix
	botapi.Bot.GetCachedItem = api.Bot.GetCachedItem
	botapi.Bot.SetCachedItems = api.Bot.SetCachedItems
	botapi.Bot.PostMessage = api.Bot.PostMessage
	botapi.Bot.Stats = api.Bot.Stats
	botapi.Bot.Build = api.Bot.Build
	botapi.UserProfile = api.UserProfile
	return
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
	// ignore spam messages
	if !tb.NoSpam(api, channel, incoming_message) {
		return
	}
	// handle commands
	if strings.HasPrefix(incoming_message, api.CommandPrefix) {
		tb.ProcessCommandMessage(api, channel, incoming_message)
	} else {
		// ignore bot messages
		if api.UserProfile.ID != "" && api.Me != "" && api.UserProfile.ID == api.Me {
			tb.logger.Println("Ignoring my own messages...")
			return
		}

		// handle text messages in separate goroutine
		go tb.processTextMessage(api, channel, incoming_message)

		// handle history (skip if sender ID is not set)
		if api.UserProfile.ID != "" && api.Me != "" {
			tb.StoreMessageHistory(api, channel, incoming_message)
		}
	}
}

func (tb *TorpedoBot) processTextMessage(api *TorpedoBotAPI, channel interface{}, incoming_message string) {
	botapi := tb.GetBotAPI(api, channel, incoming_message)
	// run message handlers here
	for idx, handler := range torpedo_registry.Config.GetTextMessageHandlers() {
		tb.logger.Printf("Running text handler #%s: %+v\n", idx, handler)
		handler(botapi, channel, incoming_message)
	}
	return
}

func (tb *TorpedoBot) Cleanup() {
	tb.logger.Printf("Running cleanup...")
	return
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
		// slow down logins
		time.Sleep(3 * time.Second)
		go wrapped(key, cmd_prefix)
	}
}

func (tb *TorpedoBot) SetBuildInfo(build, buildDate, goversion, version, projecturl string) {
	tb.Build.Build = build
	tb.Build.BuildDate = buildDate
	tb.Build.GoVersion = goversion
	tb.Build.Version = version
	tb.Build.ProjectURL = projecturl
	return
}

func (tb *TorpedoBot) RunPreParsers() {
	for pname, pfunc := range torpedo_registry.Config.GetPreParsers() {
		// TODO: Use proper logger instead
		if torpedo_registry.Config.GetConfig()["debug"] == "yes" {
			tb.logger.Printf("Running argument preparser %s\n", pname)
		}
		pfunc(torpedo_registry.Config)
	}
	return
}

func (tb *TorpedoBot) RunPostParsers() {
	for pname, pfunc := range torpedo_registry.Config.GetPostParsers() {
		// TODO: Use proper logger instead
		if torpedo_registry.Config.GetConfig()["debug"] == "yes" {
			tb.logger.Printf("Running argument postparser %s\n", pname)
		}
		pfunc(torpedo_registry.Config)
	}
	return
}

func (tb *TorpedoBot) RunCoroutines() {
	for cname, cfunc := range torpedo_registry.Config.GetCoroutines() {
		// TODO: Use proper logger instead
		if torpedo_registry.Config.GetConfig()["debug"] == "yes" {
			tb.logger.Printf("Running coroutine: %s\n", cname)
		}
		go cfunc(torpedo_registry.Config)
	}
}

func New() *TorpedoBot {
	cleanup_channel := make(chan os.Signal)
	signal.Notify(cleanup_channel, os.Interrupt, os.Kill)

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

		// Add signal listener
		go func() {
			for {
				sig := <-cleanup_channel
				bot.Cleanup()
				bot.logger.Printf("\nGot %q signal. Exiting...\n", sig)

				exitCode := 1
				if sysSig, ok := sig.(syscall.Signal); ok {
					exitCode = int(sysSig)
				}
				os.Exit(exitCode + 128)
			}
		}()

	})
	return bot
}

func (tb *TorpedoBot) RunLoop() {
	if tb.Stats.TotalAccounts > 0 {
		for {
			time.Sleep(time.Second)
		}
	} else {
		tb.logger.Fatal("No accounts configured, exiting...\n")
	}
}
