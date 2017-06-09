package multibot

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleLineMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []RichMessage) {
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

func (tb *TorpedoBot) RunLineBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1

	logger := log.New(os.Stdout, "line-bot: ", log.Lshortfile|log.LstdFlags)

	bot, err := linebot.New(strings.Split(apiKey, ":")[0],
		strings.Split(apiKey, ":")[1])
	if err != nil {
		logger.Fatal(err)
	}

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
					botApi.CommandPrefix = cmd_prefix
					go tb.processChannelEvent(botApi, channel, message.Text)
				default:
					logger.Printf("Got message type %T\n", message)

				}
			} else {
				logger.Printf("Got event type %T\n", event)
			}
		}
	})

	tb.logger.Printf("Serving Line bot on %s\n", tb.Config.LineIncomingAddr)

	if err := http.ListenAndServe(tb.Config.LineIncomingAddr, nil); err != nil {
		logger.Fatal(err)
	}
	tb.Stats.ConnectedAccounts -= 1
}
