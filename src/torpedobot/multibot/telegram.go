package multibot

import (
	"log"
	"os"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"time"
)

func (tb *TorpedoBot) RunTelegramBot(apiKey, cmd_prefix string) {
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
}
