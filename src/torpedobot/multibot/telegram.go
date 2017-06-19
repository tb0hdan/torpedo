package multibot

import (
	"os"
	"time"

	"flag"
	"regexp"
	"torpedobot/common"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
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

func (tb *TorpedoBot) ConfigureTelegramBot() {
	tb.Config.TelegramAPIKey = *flag.String("telegram", "", "Comma separated list of Telegram bot keys")
}

func (tb *TorpedoBot) RunTelegramBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1

	cu := &common.Utils{}

	if tb.Config.TelegramAPIKey == "" {
		tb.Config.TelegramAPIKey = common.GetStripEnv("TELEGRAM")
	}

	logger := cu.NewLog("telegram-bot")

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

		// handle multible bot presence
		r := regexp.MustCompile(`(?i)@(.+)bot`)
		message := r.ReplaceAllString(update.Message.Text, "")

		logger.Printf("[%s] %s\n", update.Message.From.UserName, message)

		go tb.processChannelEvent(botApi, update.Message.Chat.ID, message)

	}
	tb.Stats.ConnectedAccounts -= 1
}
