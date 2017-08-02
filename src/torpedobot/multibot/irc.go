package multibot

import (
	"flag"

	"strings"

	"crypto/tls"
	"fmt"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
	"github.com/thoj/go-ircevent"
	"gopkg.in/mgo.v2/bson"
)

var IRCAPIKey *string

func (tb *TorpedoBot) ConfigureIRCBot(cfg *torpedo_registry.ConfigStruct) {
	IRCAPIKey = flag.String("ircapikey", "", "Comma separated list of IRC creds, server:port:usessl, e.g. example.com:6679:1")
}

func (tb *TorpedoBot) ParseIRCBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("ircapikey", *JabberAPIKey)
	if cfg.GetConfig()["ircrapikey"] == "" {
		cfg.SetConfig("ircpikey", common.GetStripEnv("IRC"))
	}
}

type IRCChatroom struct {
	MyJID    string
	Chatroom string
}

type IRCAPI struct {
	ServerString string
	UseTLS       bool
}

func (ircapi *IRCAPI) Send(channel, message string, attachments ...*SkypeAttachment) {

}

func HandleIRCMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	switch api := tba.API.(type) {
	case *IRCAPI:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			api.Send(channel.(string), richmsgs[0].Text, ToSkypeAttachment(richmsgs[0]))
		} else {
			api.Send(channel.(string), message)
		}

	}
}

func (tb *TorpedoBot) RunIRCBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1
	cu := &common.Utils{}
	logger := cu.NewLog("irc-bot")
	server := strings.Split(apiKey, ":")[0]
	port := strings.Split(apiKey, ":")[1]
	usessl := strings.Split(apiKey, "@")[2]

	tb.RegisteredProtocols["*multibot.IRCAPI"] = HandleIRCMessage

	ircnick1 := "blatiblat"
	irccon := irc.IRC(ircnick1, "john doe")
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
	results := make([]*JabberChatroom, 0)
	err = collection.Find(bson.M{"myjid": GetStrippedJID(talk)}).All(&results)
	if err != nil {
		logger.Printf("No rooms available to join: %+v\n", err)
	}
	session.Close()
	for _, room := range results {
		logger.Printf("Joining IRC chatroom: %s\n", room.Chatroom)
		//talk.JoinMUCNoHistory(room.Chatroom, "TorpedoBot")
		irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(room.Chatroom) })
	}
	// end of names
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			if !strings.HasPrefix(event.Message(), "!") {
				return
			}
			//ircobj.Privmsg("<nickname | #channel>", "msg") // sends a message to either a certain nick or a channel
			fmt.Println(event.Message())    //event.Message() contains the message
			fmt.Println(event.Nick)         //event.Nick Contains the sender
			fmt.Println(event.Arguments[0]) //event.Arguments[0] Contains the channel
			if strings.HasPrefix(event.Arguments[0], "#") {
				// public msg
				irccon.Privmsg(event.Arguments[0], "lol")
			} else {
				// private msg
				irccon.Privmsg(event.Nick, "lol")
			}
		}(event)
	})
	//
	err := irccon.Connect(server)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	irccon.Loop()
	fmt.Println("connection terminated")

	// blocking run here
	tb.Stats.ConnectedAccounts -= 1
}
