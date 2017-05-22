package multibot

import (
	"os"

	"github.com/mattn/go-xmpp"
	"github.com/nlopes/slack"
	"github.com/paked/messenger"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type TorpedoBotAPI struct {
	API           interface{}
	CommandPrefix string
	Bot           *TorpedoBot
}

// https://developers.facebook.com/docs/messenger-platform/send-api-reference
const FACEBOOK_TEXT_MAX = 640

func (tba *TorpedoBotAPI) PostMessage(channel interface{}, message string, richmsgs ...RichMessage) {
	var params slack.PostMessageParameters

	switch api := tba.API.(type) {
	case *slack.Client:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
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
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
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
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			api.Send(channel.(string), richmsgs[0].Text, richmsgs[0].ToSkypeAttachment())
		} else {
			api.Send(channel.(string), message)
		}
	case *messenger.Response:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			msg, url := richmsgs[0].ToFacebookAttachment()
			if len(msg) > FACEBOOK_TEXT_MAX {
				var new_str string
				for i := 0; i < len(msg); i++ {
					if len(new_str) < FACEBOOK_TEXT_MAX {
						new_str += string(msg[i])
					} else {
						api.Text(new_str)
						new_str = ""
						new_str += string(msg[i])
					}
				}
				api.Text(new_str)
			} else {
				api.Text(msg)
			}
			api.Attachment(messenger.ImageAttachment, url)
		} else {
			api.Text(message)
		}
	}
}
