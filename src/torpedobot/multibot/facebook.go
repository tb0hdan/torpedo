package multibot

import (
	"github.com/paked/messenger"
	"strings"
	"time"
	"net/http"
	"log"
	"os"
)
func (tb *TorpedoBot) RunFacebookBot(apiKey, cmd_prefix string) {
	logger := log.New(os.Stdout, "facebook-bot: ", log.Lshortfile|log.LstdFlags)

	pageToken := strings.Split(apiKey, ":")[0]
	verifyToken := strings.Split(apiKey, ":")[1]
	client := messenger.New(messenger.Options{
		Verify:      true,
		VerifyToken: verifyToken,
		Token:       pageToken,
	})
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		logger.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		botApi := &TorpedoBotAPI{}
		botApi.API = r
		botApi.Bot = tb
		botApi.CommandPrefix = cmd_prefix

		go tb.processChannelEvent(botApi, m.Sender.ID, m.Text)
	})
	// Setup a handler to be triggered when a message is delivered
	client.HandleDelivery(func(d messenger.Delivery, r *messenger.Response) {
		logger.Println("Delivered at:", d.Watermark().Format(time.UnixDate))
	})

	// Setup a handler to be triggered when a message is read
	client.HandleRead(func(m messenger.Read, r *messenger.Response) {
		logger.Println("Read at:", m.Watermark().Format(time.UnixDate))
	})

	logger.Printf("Serving messenger bot on %s\n", tb.Config.FacebookIncomingAddr)

	if err := http.ListenAndServe(tb.Config.FacebookIncomingAddr, client.Handler()); err != nil {
		logger.Fatal(err)
	}
}
