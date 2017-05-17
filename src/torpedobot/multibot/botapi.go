package multibot

import (
	"github.com/nlopes/slack"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"github.com/mattn/go-xmpp"
	"os"
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
		if len(richmsgs) > 0 && ! richmsgs[0].IsEmpty() {
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
		var tmp string
		if len(richmsgs) > 0 && ! richmsgs[0].IsEmpty() {
			msg, tmp = richmsgs[0].ToTelegramAttachment(channel.(int64))
			api.Send(tgbotapi.NewMessage(channel.(int64), richmsgs[0].Text))
		} else {
			msg = tgbotapi.NewMessage(channel.(int64), message)
		}
		api.Send(msg)
		if tmp != "" {
			os.Remove(tmp)
		}
	case *xmpp.Client:
		msg := xmpp.Chat{}
		msg.Remote = channel.(string)
		msg.Text = message
		api.Send(msg)
	case *SkypeAPI:
		if len(richmsgs) > 0 && ! richmsgs[0].IsEmpty() {
			api.Send(channel.(string), richmsgs[0].Text, richmsgs[0].ToSkypeAttachment())
		} else {
			api.Send(channel.(string), message)
		}
	}
}

