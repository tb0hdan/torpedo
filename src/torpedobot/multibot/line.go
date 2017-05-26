package multibot

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
	"log"
	"net/http"
	"strings"
)


func (tb *TorpedoBot) RunLineBot(apiKey, cmd_prefix string) {
	logger := log.New(os.Stdout, "line-bot: ", log.Lshortfile|log.LstdFlags)

	bot, err := linebot.New(strings.Split(apiKey, ":")[0],
				strings.Split(apiKey, ":")[1])
	if err != nil {
		logger.Fatal(err)
	}

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
}

