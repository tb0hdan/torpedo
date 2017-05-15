package multibot

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-xmpp"
)

func (tb *TorpedoBot) RunJabberBot(apiKey, cmd_prefix string) {
	var talk *xmpp.Client
	var err error
	logger := log.New(os.Stdout, "jabber-bot: ", log.Lshortfile|log.LstdFlags)
	str_jid := strings.Split(apiKey, ":")[0]
	password := strings.Split(apiKey, ":")[1]
	server := strings.Split(str_jid, "@")[1]
	options := xmpp.Options{Host: server,
		User:          str_jid,
		Password:      password,
		NoTLS:         true,
		Debug:         true,
		Session:       false,
		Status:        "xa",
		StatusMessage: "",
	}

	talk, err = options.NewClient()

	if err != nil {
		logger.Fatal(err)
	}

	botApi := &TorpedoBotAPI{}
	botApi.API = talk
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	startup_ts := time.Now().Unix()
	for {
		chat, err := talk.Recv()
		if err != nil {
			log.Fatal(err)
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			passed := int64(time.Now().Unix()) - int64(startup_ts)
			logger.Println(v.Remote, v.Text, v.Stamp.Unix())
			// Since v.Stamp returns default value, use some time to catch up on messages
			if passed > 30 {
				go tb.processChannelEvent(botApi, v.Remote, v.Text)
			}
		case xmpp.Presence:
			logger.Println(v.From, v.Show)
		}
	}

}
