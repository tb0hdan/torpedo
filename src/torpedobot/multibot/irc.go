package multibot

import (
	"flag"

	"strings"

	"crypto/tls"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
	"gopkg.in/mgo.v2/bson"
	"github.com/thoj/go-ircevent"
)

var IRCAPIKey *string

func (tb *TorpedoBot) ConfigureIRCBot(cfg *torpedo_registry.ConfigStruct) {
	IRCAPIKey = flag.String("ircapikey", "",
		"Comma separated list of IRC creds, server:port:usessl, e.g. example.com:6679:1")
}

func (tb *TorpedoBot) ParseIRCBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("ircapikey", *IRCAPIKey)
	if cfg.GetConfig()["ircapikey"] == "" {
		cfg.SetConfig("ircapikey", common.GetStripEnv("IRC"))
	}
}

type IRCChatroom struct {
	MyServer string
	Channel  string
}

type IRCAPI struct {
	Connection *irc.Connection
	Event      *irc.Event
}

func (ircapi *IRCAPI) Send(channel, message string, attachments ...*SkypeAttachment) {
	if strings.HasPrefix(channel, "#") {
		// public msg
		for _, line := range strings.Split(message, "\n") {
			ircapi.Connection.Privmsg(channel, line)
		}
	} else {
		// private msg
		ircapi.Connection.Privmsg(ircapi.Event.Nick, message)
	}
}

func HandleIRCMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	switch api := tba.API.(type) {
	case *IRCAPI:
		/*if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			api.Send(channel.(string), richmsgs[0].Text, ToIRCAttachment(richmsgs[0]))
		} else { */
		api.Send(channel.(string), message)
		//}

	}
}

func (tb *TorpedoBot) RunIRCBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1
	cu := &common.Utils{}
	logger := cu.NewLog("irc-bot")
	server := strings.Split(apiKey, ":")[0]
	port := strings.Split(apiKey, ":")[1]
	usessl := strings.Split(apiKey, ":")[2]

	tb.RegisteredProtocols["*multibot.IRCAPI"] = HandleIRCMessage

	irccon := irc.IRC("torpedobot", "torpedo bot")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = usessl == "1"
	if irccon.UseTLS {
		irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	// welcome
	session, collection, err := tb.Database.GetCollection("ircChatrooms")
	if err != nil {
		logger.Fatal("Could not connect to database: %+v\n", err)
	}
	results := make([]*IRCChatroom, 0)
	err = collection.Find(bson.M{"myserver": server}).All(&results)
	if err != nil {
		logger.Printf("No rooms available to join: %+v\n", err)
	}
	session.Close()
	for _, room := range results {
		logger.Printf("Joining IRC chatroom: %s\n", room.Channel)
		irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(room.Channel) })
	}
	// end of names
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("INVITE", func(e *irc.Event) {
		session, collection, err := tb.Database.GetCollection("ircChatrooms")
		if err != nil {
			logger.Fatal("Could not connect to database: %+v\n", err)
		}
		result := IRCChatroom{}
		err = collection.Find(bson.M{"myserver": server, "channel": e.Arguments[1]}).One(&result)
		if err != nil {
			logger.Println(err)
			// no record, insert new one
			err = collection.Insert(&IRCChatroom{MyServer: server, Channel: e.Arguments[1]})
			if err != nil {
				logger.Fatal(err)
			}
			// join new room
			irccon.Join(e.Arguments[1])
		}
		session.Close()
	})
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			botApi := &TorpedoBotAPI{}
			api := &IRCAPI{Connection: irccon, Event: event}
			botApi.API = api
			botApi.Bot = tb
			botApi.CommandPrefix = cmd_prefix
			tb.processChannelEvent(botApi, event.Arguments[0], event.Message())
		}(event)
	})
	//
	err = irccon.Connect(server + ":" + port)
	if err != nil {
		logger.Printf("Err %s", err)
		return
	}
	// blocking run here
	irccon.Loop()

	// we'll probably won't get here
	logger.Println("connection terminated")
	tb.Stats.ConnectedAccounts -= 1
}
