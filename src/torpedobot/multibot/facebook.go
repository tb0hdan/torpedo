package multibot

import (
	"net/http"
	"strings"
	"time"

	"flag"
	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"

	"github.com/paked/messenger"
)

// https://developers.facebook.com/docs/messenger-platform/send-api-reference
const FACEBOOK_TEXT_MAX = 640
var (
	FacebookAPIKey *string
	FacebookIncomingAddr *string
)

func HandleFacebookMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	switch api := tba.API.(type) {
	case *messenger.Response:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			msg, url := richmsgs[0].ToGenericAttachment()
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

func (tb *TorpedoBot) ConfigureFacebookBot(cfg *torpedo_registry.ConfigStruct) {
	FacebookAPIKey = flag.String("facebook", "", "Comma separated list of Facebook creds, page_token1:verify_token1,..")
	FacebookIncomingAddr = flag.String("facebook_incoming_addr", "0.0.0.0:3979", "Listen on this address for incoming Facebook messages")

}

func (tb *TorpedoBot) ParseFacebookBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("facebookapikey", *FacebookAPIKey)
	cfg.SetConfig("facebookincomingaddr", *FacebookIncomingAddr)
	if cfg.GetConfig()["facebookapikey"] == "" {
		cfg.SetConfig("facebookapikey", common.GetStripEnv("FACEBOOK"))
	}
}

func (tb *TorpedoBot) RunFacebookBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1
	cu := &common.Utils{}
	logger := cu.NewLog("facebook-bot")

	tb.RegisteredProtocols["*messenger.Response"] = HandleFacebookMessage

	pageToken := strings.Split(apiKey, ":")[0]
	verifyToken := strings.Split(apiKey, ":")[1]
	appSecret := strings.Split(apiKey, ":")[2]
	client := messenger.New(messenger.Options{
		AppSecret:   appSecret,
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
		botApi.UserProfile = &torpedo_registry.UserProfile{}

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

	logger.Printf("Serving messenger bot on %s\n", torpedo_registry.Config.GetConfig()["facebookincomingaddr"])

	if err := http.ListenAndServe(torpedo_registry.Config.GetConfig()["facebookincomingaddr"], client.Handler()); err != nil {
		logger.Fatal(err)
	}
	tb.Stats.ConnectedAccounts -= 1
}
