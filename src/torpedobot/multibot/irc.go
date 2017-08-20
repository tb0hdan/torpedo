package multibot

import (
	"flag"

	"strings"

	"crypto/tls"

	"fmt"

	"log"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
	irc "github.com/thoj/go-ircevent"
	"gopkg.in/mgo.v2/bson"
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
			if strings.TrimSpace(line) == "" {
				continue
			}
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
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			msg, url := richmsgs[0].ToGenericAttachment()
			api.Send(channel.(string), fmt.Sprintf("%s\n%s", msg, url))
		} else {
			api.Send(channel.(string), message)
		}

	}
}

// set custom logger + version
func (tb *TorpedoBot) myIRC(nick, user string, log *log.Logger) *irc.Connection {
	connection := irc.IRC(nick, user)
	connection.Log = log
	connection.Version = fmt.Sprintf("%s v%s (%s)", tb.Build.ProjectURL, tb.Build.Version, tb.Build.Build)
	return connection
}

func (tb *TorpedoBot) RunIRCBot(apiKey, cmd_prefix string) {
	var (
		nick, server string
	)
	tb.Stats.ConnectedAccounts += 1
	cu := &common.Utils{}
	logger := cu.NewLog("irc-bot")
	tb.RegisteredProtocols["*multibot.IRCAPI"] = HandleIRCMessage

	user_server := strings.Split(apiKey, ":")[0]
	if len(strings.Split(user_server, "@")) == 2 {
		nick = strings.Split(user_server, "@")[0]
		server = strings.Split(user_server, "@")[1]
	} else {
		nick = "torpedobot"
		server = user_server
	}
	port := strings.Split(apiKey, ":")[1]
	usessl := strings.Split(apiKey, ":")[2]

	irccon := tb.myIRC(nick, fmt.Sprintf("%s bot", nick), logger)
	// TODO: Add -v switch to enable these
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true

	// TLS Config
	irccon.UseTLS = usessl == "1"
	if irccon.UseTLS {
		irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Password config
	if len(strings.Split(apiKey, ":")) > 3 {
		irccon.Password = strings.Split(apiKey, ":")[3]
	}

	//welcome
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
			botApi.UserProfile = &torpedo_registry.UserProfile{ID: event.User, Nick: event.Nick}
			botApi.Me = irccon.GetNick()

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
