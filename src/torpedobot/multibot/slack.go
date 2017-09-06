package multibot

import (
	"strconv"
	"time"

	common "github.com/tb0hdan/torpedo_common"

	"flag"

	"github.com/nlopes/slack"
	"github.com/tb0hdan/torpedo_registry"
)

var SlackAPIKey *string

func ToSlackAttachment(rm torpedo_registry.RichMessage) (params slack.PostMessageParameters) {
	attachment := slack.Attachment{
		Color:     rm.BarColor,
		Text:      rm.Text,
		Title:     rm.Title,
		TitleLink: rm.TitleLink,
		ImageURL:  rm.ImageURL,
	}
	params.Attachments = []slack.Attachment{attachment}
	return
}

func HandleSlackMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	var params slack.PostMessageParameters
	if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
		params = ToSlackAttachment(richmsgs[0])
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

func (tb *TorpedoBot) ConfigureSlackBot(cfg *torpedo_registry.ConfigStruct) {
	SlackAPIKey = flag.String("slack", "", "Comma separated list of Slack legacy tokens")

}

func (tb *TorpedoBot) ParseSlackBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("slackapikey", *SlackAPIKey)
	if cfg.GetConfig()["slackapikey"] == "" {
		cfg.SetConfig("slackapikey", common.GetStripEnv("SLACK"))
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
	botApi.UserProfile = &torpedo_registry.UserProfile{}

	tb.RegisteredProtocols["*slack.Client"] = HandleSlackMessage

	for msg := range rtm.IncomingEvents {
		logger.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			logger.Println("Infos:", ev.Info)
			logger.Println("Connection counter:", ev.ConnectionCount)
			botApi.Me = ev.Info.User.ID
			// Replace #general with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			logger.Printf("Message: %v\n", ev)
			if ev.Type == "message" && ev.User != "" {
				user, err := api.GetUserInfo(ev.User)
				if err == nil {
					botApi.UserProfile = &torpedo_registry.UserProfile{Nick: user.Name,
						RealName: user.RealName,
						Timezone: user.TZ,
						Phone:    user.Profile.Phone,
						Email:    user.Profile.Email,
						IsBot:    user.IsBot,
						ID:       user.ID,
					}
				} else {
					logger.Printf("Error getting user info for %s\n", ev.User)
				}

				channel := ev.Channel
				incoming_message := ev.Text
				messageTS, _ := strconv.ParseFloat(ev.Timestamp, 64)
				jitter := int64(time.Now().Unix()) - int64(messageTS)
				if jitter < 10 {
					go tb.processChannelEvent(botApi, channel, incoming_message)
				}
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
