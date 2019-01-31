package multibot

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	common "github.com/tb0hdan/torpedo_common"

	"github.com/nlopes/slack"
	"github.com/tb0hdan/torpedo_registry"
)

var SlackAPIKey *string

func ToSlackAttachment(rm torpedo_registry.RichMessage) (msg slack.MsgOption) {

	attachment := slack.Attachment{
		Color:     rm.BarColor,
		Text:      rm.Text,
		Title:     rm.Title,
		TitleLink: rm.TitleLink,
		ImageURL:  rm.ImageURL,
	}
	msg = slack.MsgOptionAttachments(attachment)
	return
}

func HandleSlackMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	var params slack.MsgOption //slack.PostMessageParameters
	if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
		params = ToSlackAttachment(richmsgs[0])
	}
	/*
	params.UnfurlLinks = true
	params.UnfurlMedia = true
	*/

	switch api := tba.API.(type) {
	case *slack.Client:

		channelID, timestamp, err := api.PostMessage(channel.(string), slack.MsgOptionText(message, false), params)
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
	account := &torpedo_registry.Account{
		APIKey:        apiKey,
		CommandPrefix: cmd_prefix,
	}
	torpedo_registry.Accounts.AppendAccounts(account)
	tb.RunSlackBotAccount(account)
}

func (tb *TorpedoBot) RunSlackBotAccount(account *torpedo_registry.Account) {
	tb.Stats.ConnectedAccounts += 1

	slack_debug := torpedo_registry.Config.GetConfig()["debug"] == "yes"

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)

	api := slack.New(account.APIKey, slack.OptionDebug(slack_debug), slack.OptionLog(logger))
	account.API = api

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	botApi := &TorpedoBotAPI{}
	botApi.API = api
	botApi.Bot = tb
	botApi.CommandPrefix = account.CommandPrefix
	botApi.UserProfile = &torpedo_registry.UserProfile{}

	tb.RegisteredProtocols["*slack.Client"] = HandleSlackMessage

	for msg := range rtm.IncomingEvents {
		// TODO: Use proper logger instead
		if torpedo_registry.Config.GetConfig()["debug"] == "yes" {
			logger.Print("Event Received: ")
		}
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			account.Connection.Connected = true
			account.Connection.ReconnectCount += 1
			// TODO: Use proper logger instead
			if torpedo_registry.Config.GetConfig()["debug"] == "yes" {
				logger.Println("Infos:", ev.Info)
				logger.Println("Connection counter:", ev.ConnectionCount)
			}
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
				// System notifications, like "you've been invited / kicked" come from USLACKBOT, ignore them...
				if jitter < 10 && botApi.UserProfile.ID != "USLACKBOT" {
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
			account.Connection.Connected = false
			logger.Printf("Invalid credentials")
			break

		default:
			// Ignore other events..
			//logger.Printf("Unexpected: %v\n", msg.Data)
		}
	}
	tb.Stats.ConnectedAccounts -= 1
}
