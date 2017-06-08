package multibot

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-xmpp"
)

func (tb *TorpedoBot) JabberServerInfo(jid, server string, c *xmpp.Client) (string, error) {
	const namespace = "http://jabber.org/protocol/disco#info"
	// use getCookie for a pseudo random id.
	reqID := strconv.FormatUint(uint64(time.Now().Unix()), 10)
	return c.RawInformationQuery(jid, server, reqID, xmpp.IQTypeGet, namespace, "")
}

func (tb *TorpedoBot) SendJabberDisco(jid, server string, client *xmpp.Client) {
	_, err := tb.JabberServerInfo(jid, server, client)
	if err != nil {
		tb.logger.Printf("%+v\n", err)
	}
	return
}

func (tb *TorpedoBot) WaitAndSendJabberDisco(jid, server string, client *xmpp.Client) {
	// Wait for event loop to start
	time.Sleep(10 * time.Second)
	tb.SendJabberDisco(jid, server, client)
	return
}

func (tb *TorpedoBot) JabberPinger(jid, server string, client *xmpp.Client) {
	sleep := 60
	tb.logger.Printf("Sending ping every %d seconds\n", sleep)
	for {
		client.PingC2S(jid, server)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func HandleJabberMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []RichMessage) {
	switch api := tba.API.(type) {
	case *xmpp.Client:
		msg := xmpp.Chat{}
		msg.Remote = channel.(string)
		msg.Text = message
		api.Send(msg)
	}
}

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

	tb.RegisteredProtocols["*xmpp.Client"] = HandleJabberMessage

	if err != nil {
		logger.Fatal(err)
	}

	botApi := &TorpedoBotAPI{}
	botApi.API = talk
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	startup_ts := time.Now().Unix()
	go tb.WaitAndSendJabberDisco(str_jid, server, talk)
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
		case xmpp.IQ:
			if v.Type == "result" && v.ID == "c2s1" {
				logger.Printf("Got pong from %s to %s\n", v.From, v.To)
			}
			if strings.Contains(string(v.Query), "urn:xmpp:ping") {
				go tb.JabberPinger(str_jid, server, talk)
			}
		default:
			logger.Printf("Unknown event: %T\n", v)
		}
	}

}
