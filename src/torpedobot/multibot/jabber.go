package multibot

import (
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-xmpp"
	"gopkg.in/mgo.v2/bson"
	"torpedobot/common"
)

type JabberChatroom struct {
	MyJID string
	Chatroom string
}

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
		msg.Type = tba.Type
		msg.Text = message
		if tba.Type == "groupchat" {
			msg.Remote = strings.Split(msg.Remote, "/")[0]
		}
		api.Send(msg)
	}
}

func GetStrippedJID(cli *xmpp.Client) (jid string) {
	jid = strings.Split(cli.JID(), "/")[0]
	return
}


func (tb *TorpedoBot) RunJabberBot(apiKey, cmd_prefix string) {
	var talk *xmpp.Client
	var err error
	tb.Stats.ConnectedAccounts += 1
	cu := &common.Utils{}
	logger := cu.NewLog("jabber-bot")
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
	// join rooms
	session, collection, err := tb.Database.GetCollection("jabberChatrooms")
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
		logger.Printf("Joining chatroom: %s\n", room.Chatroom)
		talk.JoinMUCNoHistory(room.Chatroom, "TorpedoBot")
	}
	//
	for {
		chat, err := talk.Recv()
		if err != nil {
			logger.Fatal(err)
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			passed := int64(time.Now().Unix()) - int64(startup_ts)
			logger.Println(v.Remote, v.Text, v.Stamp.Unix(), v.Type, v.Other, v.OtherElem)
			for _, element := range v.OtherElem {
				if element.XMLName.Space == "jabber:x:conference" {
					session, collection, err := tb.Database.GetCollection("jabberChatrooms")
					if err != nil {
						logger.Fatal("Could not connect to database: %+v\n", err)
					}
					result := JabberChatroom{}
					err = collection.Find(bson.M{"myjid": GetStrippedJID(talk), "chatroom": v.Remote}).One(&result)
					if err != nil {
						logger.Println(err)
						// no record, insert new one
						err = collection.Insert(&JabberChatroom{GetStrippedJID(talk), v.Remote})
						if err != nil {
							logger.Fatal(err)
						}
						// join new room
						talk.JoinMUCNoHistory(v.Remote, "TorpedoBot")
					}
					session.Close()
					break
				}
			}
			// Since v.Stamp returns default value, use some time to catch up on messages
			if passed > 30 {
				botApi.Type = v.Type
				go tb.processChannelEvent(botApi, v.Remote, v.Text)
			}
		case xmpp.Presence:
			if v.Type == "subscribe" {
				talk.ApproveSubscription(v.From)
			}
			logger.Println(v.From, v.Show, v.Type)
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
	tb.Stats.ConnectedAccounts -= 1

}
