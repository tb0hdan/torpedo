package multibot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"torpedobot/common"
	"torpedobot/memcache"



	"github.com/mattn/go-xmpp"
	"github.com/nlopes/slack"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)


type TorpedoBot struct {
	caches          map[string]*memcache.MemCacheType
	commandHandlers map[string]func(*TorpedoBotAPI, interface{}, string)
	config          struct {
	}
}


type TorpedoBotAPI struct {
	API        interface{}
	CommandPrefix string
	Bot *TorpedoBot
}


func (tba *TorpedoBotAPI) PostMessage(channel interface{}, message string, parameters ...interface{}) {
	var params slack.PostMessageParameters

	switch api := tba.API.(type) {
	case *slack.Client:
		if len(parameters) > 0 {
			params = parameters[0].(slack.PostMessageParameters)
		}
		channelID, timestamp, err := api.PostMessage(channel.(string), message, params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	case *tgbotapi.BotAPI:
		msg := tgbotapi.NewMessage(channel.(int64), message)
		api.Send(msg)
	case *xmpp.Client:
		msg := xmpp.Chat{}
		msg.Remote = channel.(string)
		msg.Text = message
		api.Send(msg)
	}
}


func (tb *TorpedoBot) PostMessage(channel interface{}, message string, api *TorpedoBotAPI, parameters ...slack.PostMessageParameters) {
	var params slack.PostMessageParameters

	if len(parameters) > 0 {
		params = parameters[0]
	}
	api.PostMessage(channel, message, params)
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
		fmt.Printf("PROCESS! -> `%s`", command)
		if found == 0 {
			api.PostMessage(channel, fmt.Sprintf("Could not process your message: %s%s. Command unknown. Send %shelp for list of valid commands.", api.CommandPrefix, command, api.CommandPrefix), api)
		}
	}
}


func (tb *TorpedoBot) RunSlackBot(apiKey, cmd_prefix string) {
	api := slack.New(apiKey)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	botApi := &TorpedoBotAPI{}
	botApi.API = api
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	// TODO: Move this somewhere else
	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			channel := ev.Channel
			incoming_message := ev.Text
			messageTS, _ := strconv.ParseFloat(ev.Timestamp, 64)
			jitter := int64(time.Now().Unix()) - int64(messageTS)
			if jitter < 10 {
				go tb.processChannelEvent(botApi, channel, incoming_message)
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			//fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

}

func (tb *TorpedoBot) RunTelegramBot(apiKey, cmd_prefix string) {
	api, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)

	botApi := &TorpedoBotAPI{}
	botApi.API = api
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	for update := range updates {
		if update.Message == nil {
			continue
		}

		jitter := int64(time.Now().Unix()) - int64(update.Message.Date)

		if jitter > 10 {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		go tb.processChannelEvent(botApi, update.Message.Chat.ID, update.Message.Text)

	}
}

func (tb *TorpedoBot) RunJabberBot(apiKey, cmd_prefix string) {
	var talk *xmpp.Client
	var err error
	str_jid := strings.Split(apiKey, ":")[0]
	password := strings.Split(apiKey, ":")[1]
	server := strings.Split(str_jid, "@")[1]
	options := xmpp.Options{Host: server,
		User:          str_jid,
		Password:      password,
		NoTLS:         true,
		Debug:         true,
		Session:       false,
		Status:        "xa",
		StatusMessage: "",
	}

	talk, err = options.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	botApi := &TorpedoBotAPI{}
	botApi.API = talk
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	startup_ts := time.Now().Unix()
	for {
		chat, err := talk.Recv()
		if err != nil {
			log.Fatal(err)
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			passed := int64(time.Now().Unix()) - int64(startup_ts)
			fmt.Println(v.Remote, v.Text, v.Stamp.Unix())
			// Since v.Stamp returns default value, use some time to catch up on messages
			if passed > 30 {
				go tb.processChannelEvent(botApi, v.Remote, v.Text)
			}
		case xmpp.Presence:
			fmt.Println(v.From, v.Show)
		}
	}

}

func (tb *TorpedoBot) RunLoop() {
	for {
		time.Sleep(time.Second)
	}
}

func (tb *TorpedoBot) RunBotsCSV(method func(apiKey string, cmd_prefix string), CSV, cmd_prefix string) {
	for _, key := range strings.Split(CSV, ",") {
		if key == "" {
			continue
		}
		go method(key, cmd_prefix)
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
		fmt.Printf("\nUsing cached quote...%v\n", cache.Len())
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

func New() (bot *TorpedoBot) {
	bot = &TorpedoBot{}
	bot.caches = make(map[string]*memcache.MemCacheType)
	return
}
