package multibot

import (
	"github.com/paked/messenger"
	"strings"
	"time"
	"net/http"
)
func (tb *TorpedoBot) RunFacebookBot(apiKey, cmd_prefix string) {
	pageToken := strings.Split(apiKey, ":")[0]
	verifyToken := strings.Split(apiKey, ":")[1]
	client := messenger.New(messenger.Options{
		Verify:      true,
		VerifyToken: verifyToken,
		Token:       pageToken,
	})
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		tb.logger.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		/*
		p, err := client.ProfileByID(m.Sender.ID)
		if err != nil {
			tb.logger.Println("Something went wrong!", err)
		}
		*/
		botApi := &TorpedoBotAPI{}
		botApi.API = r
		botApi.Bot = tb
		botApi.CommandPrefix = cmd_prefix

		//r.Text(fmt.Sprintf("Hello, %v!", p.FirstName))
		go tb.processChannelEvent(botApi, m.Sender.ID, m.Text)
	})
	// Setup a handler to be triggered when a message is delivered
	client.HandleDelivery(func(d messenger.Delivery, r *messenger.Response) {
		tb.logger.Println("Delivered at:", d.Watermark().Format(time.UnixDate))
	})

	// Setup a handler to be triggered when a message is read
	client.HandleRead(func(m messenger.Read, r *messenger.Response) {
		tb.logger.Println("Read at:", m.Watermark().Format(time.UnixDate))
	})

	tb.logger.Printf("Serving messenger bot on %s\n", tb.config.FacebookIncomingAddr)

	http.ListenAndServe(tb.config.FacebookIncomingAddr, client.Handler())
}
