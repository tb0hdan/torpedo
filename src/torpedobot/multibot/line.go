package multibot

import (
	"flag"
	"net/http"
	"strings"

	common "github.com/tb0hdan/torpedo_common"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tb0hdan/torpedo_registry"
)

var (
	LineAPIKey       *string
	LineIncomingAddr *string
)

func HandleLineMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	switch api := tba.API.(type) {
	case *linebot.Client:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			msg, url := richmsgs[0].ToGenericAttachment()
			// Use replyToken as channel
			api.PushMessage(channel.(string), linebot.NewTextMessage(msg)).Do()
			api.PushMessage(channel.(string), linebot.NewImageMessage(url, url)).Do()

		} else {
			// Use replyToken as channel
			api.PushMessage(channel.(string), linebot.NewTextMessage(message)).Do()
		}
	}
}

func (tb *TorpedoBot) ConfigureLineBot(cfg *torpedo_registry.ConfigStruct) {
	LineAPIKey = flag.String("line", "", "Line.Me credentials client_secret:client_token,")
	LineIncomingAddr = flag.String("line_incoming_addr", "0.0.0.0:3981", "Listen on this address for incoming Line.Me messages")

}

func (tb *TorpedoBot) ParseLineBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("lineapikey", *LineAPIKey)
	cfg.SetConfig("lineincomingaddr", *LineIncomingAddr)
	if cfg.GetConfig()["lineapikey"] == "" {
		cfg.SetConfig("lineapikey", common.GetStripEnv("LINE"))
	}
}

func (tb *TorpedoBot) RunLineBot(apiKey, cmd_prefix string) {
	account := &torpedo_registry.Account{
		APIKey:        apiKey,
		CommandPrefix: cmd_prefix,
	}
	torpedo_registry.Accounts.AppendAccounts(account)
	tb.RunLineBotAccount(account)
}

func (tb *TorpedoBot) RunLineBotAccount(account *torpedo_registry.Account) {
	tb.Stats.ConnectedAccounts += 1
	account.Connection.ReconnectCount += 1

	cu := &common.Utils{}

	logger := cu.NewLog("line-bot")

	bot, err := linebot.New(strings.Split(account.APIKey, ":")[0],
		strings.Split(account.APIKey, ":")[1])
	if err != nil {
		logger.Fatal(err)
	}

	account.API = bot
	tb.RegisteredProtocols["*linebot.Client"] = HandleLineMessage

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				var channel string
				if event.Source.GroupID != "" {
					channel = event.Source.GroupID
				} else if event.Source.RoomID != "" {
					channel = event.Source.RoomID
				} else if event.Source.UserID != "" {
					channel = event.Source.UserID
				}
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					botApi := &TorpedoBotAPI{}
					botApi.API = bot
					botApi.Bot = tb
					botApi.CommandPrefix = account.CommandPrefix
					botApi.UserProfile = &torpedo_registry.UserProfile{ID: channel}
					botApi.Me = "torpedobot"

					go tb.processChannelEvent(botApi, channel, message.Text)
				default:
					logger.Printf("Got message type %T\n", message)

				}
			} else {
				logger.Printf("Got event type %T\n", event)
			}
		}
	})

	tb.logger.Printf("Serving Line bot on %s\n", torpedo_registry.Config.GetConfig()["lineincomingaddr"])

	if err := http.ListenAndServe(torpedo_registry.Config.GetConfig()["lineincomingaddr"], nil); err != nil {
		logger.Fatal(err)
	}
	tb.Stats.ConnectedAccounts -= 1
}
