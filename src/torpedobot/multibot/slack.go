package multibot

import (
	"strconv"
	"time"

	"torpedobot/common"

	"github.com/nlopes/slack"
)

func HandleSlackMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []RichMessage) {
	var params slack.PostMessageParameters
	if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
		params = richmsgs[0].ToSlackAttachment()
	}

	switch api := tba.API.(type) {
	case *slack.Client:
		channelID, timestamp, err := api.PostMessage(channel.(string), message, params)
		if err != nil {
			tba.Bot.logger.Printf("%s\n", err)
			return
		}
		tba.Bot.logger.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
}

func (tb *TorpedoBot) RunSlackBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1

	api := slack.New(apiKey)
	cu := &common.Utils{}
	logger := cu.NewLog("slack-bot")
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	botApi := &TorpedoBotAPI{}
	botApi.API = api
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	tb.RegisteredProtocols["*slack.Client"] = HandleSlackMessage

	for msg := range rtm.IncomingEvents {
		logger.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			logger.Println("Infos:", ev.Info)
			logger.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			logger.Printf("Message: %v\n", ev)
			channel := ev.Channel
			incoming_message := ev.Text
			messageTS, _ := strconv.ParseFloat(ev.Timestamp, 64)
			jitter := int64(time.Now().Unix()) - int64(messageTS)
			if jitter < 10 {
				go tb.processChannelEvent(botApi, channel, incoming_message)
			}

		case *slack.PresenceChangeEvent:
			logger.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			logger.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			logger.Printf("Invalid credentials")
			break

		default:
			// Ignore other events..
			//logger.Printf("Unexpected: %v\n", msg.Data)
		}
	}
	tb.Stats.ConnectedAccounts -= 1
}
