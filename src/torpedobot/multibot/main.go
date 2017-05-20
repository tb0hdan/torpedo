package multibot


import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"torpedobot/common"
	"torpedobot/memcache"
	"sync"

	"github.com/getsentry/raven-go"
)

var bot *TorpedoBot
var once sync.Once

type TorpedoBot struct {
	caches          map[string]*memcache.MemCacheType
	commandHandlers map[string]func(*TorpedoBotAPI, interface{}, string)
	config          struct {
		FacebookIncomingAddr string
		SkypeIncomingAddr string
	}
	logger *log.Logger
}

func (tb *TorpedoBot) PostMessage(channel interface{}, message string, api *TorpedoBotAPI, richmsgs ...RichMessage) {
	if len(richmsgs) > 0 {
		api.PostMessage(channel, message, richmsgs[0])
	} else {
		api.PostMessage(channel, message)
	}

}

func (tb *TorpedoBot) processChannelEvent(api *TorpedoBotAPI, channel interface{}, incoming_message string) {
	if strings.HasPrefix(incoming_message, api.CommandPrefix) {
		command := strings.TrimPrefix(incoming_message, api.CommandPrefix)
		found := 0
		for handler := range tb.commandHandlers {
			if strings.HasPrefix(strings.Split(command, " ")[0], handler) {
				found += 1
				tb.commandHandlers[handler](api, channel, incoming_message)
				break
			}
		}
		tb.logger.Printf("PROCESS! -> `%s`", command)
		if found == 0 {
			api.PostMessage(channel, fmt.Sprintf("Could not process your message: %s%s. Command unknown. Send %shelp for list of valid commands.", api.CommandPrefix, command, api.CommandPrefix))
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
	env_dsn := os.Getenv("SENTRY_DSN")
	if  env_dsn != "" {
		tb.logger.Print("Using Sentry error reporting...\n")
		raven.SetDSN(env_dsn)
		wrapped = func(apiKey, cmd_prefix string) {
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
		go wrapped(key, cmd_prefix)
	}
}

func (tb *TorpedoBot) RegisterHandlers(handlers map[string]func(*TorpedoBotAPI, interface{}, string)) {
	tb.commandHandlers = handlers
	return
}

func (tb *TorpedoBot) GetCommandHandlers() (handlers map[string]func(*TorpedoBotAPI, interface{}, string)) {
	return tb.commandHandlers
}

func (tb *TorpedoBot) GetCreateCache(name string) (cache *memcache.MemCacheType) {
	value, success := tb.caches[name]
	if !success {
		cache = memcache.New()
		tb.caches[name] = cache
	} else {
		cache = value
	}
	return
}

func (tb *TorpedoBot) GetCachedItem(name string) (item string) {
	cache := *tb.GetCreateCache(name)
	if cache.Len() > 0 {
		tb.logger.Printf("\nUsing cached quote...%v\n", cache.Len())
		key := ""
		for key = range cache.Cache() {
			break
		}
		quote, _ := cache.Get(key)
		cache.Delete(key)
		item = quote
	}
	return
}

func (tb *TorpedoBot) SetCachedItems(name string, items map[int]string) (item string) {
	cache := *tb.GetCreateCache(name)
	for idx := range items {
		message := common.MD5Hash(items[idx])
		_, ok := cache.Get(message)
		if !ok {
			cache.Set(message, items[idx])
		}
	}

	item = items[0]
	message := common.MD5Hash(item)
	cache.Delete(message)
	return
}


func New(facebook_incoming_addr, skype_incoming_addr string) *TorpedoBot {
	once.Do(func() {
		bot = &TorpedoBot{}
		bot.logger = log.New(os.Stdout, "torpedo-bot: ", log.Lshortfile|log.LstdFlags)
		bot.caches = make(map[string]*memcache.MemCacheType)
		bot.config.SkypeIncomingAddr = skype_incoming_addr
		bot.config.FacebookIncomingAddr = facebook_incoming_addr
	})
	return bot
}
