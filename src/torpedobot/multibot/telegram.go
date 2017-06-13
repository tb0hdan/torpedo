package multibot

import (
	"log"
	"os"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"torpedobot/common"
)

func (rm *RichMessage) ToTelegramAttachment(channel int64) (msg tgbotapi.Chattable, fname string) {
	cu := &common.Utils{}
	fname, _, is_image, err := cu.DownloadToTmp(rm.ImageURL)
	if is_image && err == nil {
		msg = tgbotapi.NewPhotoUpload(channel, fname)
	}
	return
}

func HandleTelegramMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []RichMessage) {
	switch api := tba.API.(type) {
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
	}
}

func (tb *TorpedoBot) RunTelegramBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1

	logger := log.New(os.Stdout, "telegram-bot: ", log.Lshortfile|log.LstdFlags)

	api, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		logger.Panic(err)
	}

	api.Debug = true

	logger.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)

	botApi := &TorpedoBotAPI{}
	botApi.API = api
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	tb.RegisteredProtocols["*tgbotapi.BotAPI"] = HandleTelegramMessage

	for update := range updates {
		if update.Message == nil {
			continue
		}

		jitter := int64(time.Now().Unix()) - int64(update.Message.Date)

		if jitter > 10 {
			continue
		}

		logger.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		go tb.processChannelEvent(botApi, update.Message.Chat.ID, update.Message.Text)

	}
	tb.Stats.ConnectedAccounts -= 1
}
