package multibot

import (
	"github.com/nlopes/slack"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"github.com/mattn/go-xmpp"
)

type TorpedoBotAPI struct {
	API           interface{}
	CommandPrefix string
	Bot           *TorpedoBot
}


func (tba *TorpedoBotAPI) PostMessage(channel interface{}, message string, richmsgs ...RichMessage) {
	var params slack.PostMessageParameters

	switch api := tba.API.(type) {
	case *slack.Client:
		if len(richmsgs) > 0 {
			params = richmsgs[0].ToSlackAttachment()
		}
		channelID, timestamp, err := api.PostMessage(channel.(string), message, params)
		if err != nil {
			tba.Bot.logger.Printf("%s\n", err)
			return
		}
		tba.Bot.logger.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	case *tgbotapi.BotAPI:
		var msg tgbotapi.Chattable
		if len(richmsgs) > 0 {
			msg = richmsgs[0].ToTelegramAttachment(channel.(int64))
		} else {
			msg = tgbotapi.NewMessage(channel.(int64), message)
		}
		api.Send(msg)
	case *xmpp.Client:
		msg := xmpp.Chat{}
		msg.Remote = channel.(string)
		msg.Text = message
		api.Send(msg)
	case *SkypeAPI:
		api.Send(channel.(string), message)
	}
}

