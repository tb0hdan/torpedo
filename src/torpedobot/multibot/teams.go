package multibot

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	common "github.com/tb0hdan/torpedo_common"
	memcache "github.com/tb0hdan/torpedo_common/memcache"
	"github.com/tb0hdan/torpedo_registry"
)

// Custom Bot is expected to reply within 5 seconds
// https://msdn.microsoft.com/en-us/microsoft-teams/custombot#sending-a-reply
const sleepMax = 5

var (
	TeamsIncomingAddr *string
	TeamsAPIKey       *string
	TeamsMessageQueue = memcache.New()
)

type TeamsAPI struct {
	GUID   string
	logger *log.Logger
}

func (sapi *TeamsAPI) Send(channel, message string, attachments ...*SkypeAttachment) {
	outgoing_message := &SkypeOutgoingMessage{Text: message,
		Type: "message",
		TextFormat:  "plain",
		Attachments: attachments}
	// Append message to TeamsMessageQueue
	if _, ok := TeamsMessageQueue.Get(sapi.GUID); !ok {
		body, _ := json.Marshal(outgoing_message)
		TeamsMessageQueue.Set(sapi.GUID, []string{string(body)})
	}
	return
}

func HandleTeamsMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	switch api := tba.API.(type) {
	case *TeamsAPI:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			api.Send(channel.(string), richmsgs[0].Text, ToSkypeAttachment(richmsgs[0]))
		} else {
			api.Send(channel.(string), message)
		}

	}
}

func (tb *TorpedoBot) ConfigureTeamsBot(cfg *torpedo_registry.ConfigStruct) {
	TeamsIncomingAddr = flag.String("teams_incoming_addr", "0.0.0.0:3982", "Listen on this address for incoming Teams messages")
	TeamsAPIKey = flag.String("teams", "", "Comma separated list of Microsoft Teams user bot secrets")
}

func (tb *TorpedoBot) ParseTeamsBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("teamsincomingaddr", *TeamsIncomingAddr)
	cfg.SetConfig("teamsapikey", *TeamsAPIKey)
	if cfg.GetConfig()["teamsapikey"] == "" {
		cfg.SetConfig("teamsapikey", common.GetStripEnv("MSTEAMS"))
	}
}

func (tb *TorpedoBot) RunTeamsBot(apiKey, cmd_prefix string) {
	tb.Stats.ConnectedAccounts += 1
	cu := &common.Utils{}
	logger := cu.NewLog("teams-bot")

	teams_api := &TeamsAPI{}
	teams_api.logger = logger

	tb.RegisteredProtocols["*multibot.TeamsAPI"] = HandleTeamsMessage

	http.HandleFunc("/api/teams-messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		defer r.Body.Close()
		body_bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			tb.logger.Fatalf("readAll failed with %+v\n", err)
			return
		}
		logger.Printf("Teams incoming message: %s\n", string(body_bytes))
		// SkypeIncomming message seems to be compatible with Teams User Bot incoming message:
		// https://msdn.microsoft.com/en-us/microsoft-teams/botsconversation#receiving-messages
		message := &SkypeIncomingMessage{}
		err = json.Unmarshal(body_bytes, message)
		if err != nil {
			logger.Fatalf("JSON unmarshalling failed with %+v\n", err)
			return
		}

		botApi := &TorpedoBotAPI{}
		teams_api.GUID = uuid.New().String()
		botApi.API = teams_api
		botApi.Bot = tb
		botApi.CommandPrefix = cmd_prefix
		botApi.UserProfile = &torpedo_registry.UserProfile{ID: message.From.ID, Nick: message.From.Name}
		// FIXME: Remove hardcode
		botApi.Me = "torpedobot"

		re := regexp.MustCompile(`^(<at>.+</at>(\s|&nbsp;))?`)
		msg := re.ReplaceAllString(message.Text, "")
		logger.Printf("Message: `%s`\n", msg)
		go tb.processChannelEvent(botApi, message.Conversation.ID, msg)

		stopFlag := false
		ticker := time.NewTicker(time.Millisecond * 100)
		go func() {
			for range ticker.C {
				body, ok := TeamsMessageQueue.Get(teams_api.GUID)
				if ok {
					w.Write([]byte(body[0]))
					TeamsMessageQueue.Delete(teams_api.GUID)
					stopFlag = true
					break
				}
			}
		}()
		for i := 0; i <= sleepMax * 10; i++ {
			if stopFlag {
				break
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		ticker.Stop()
		logger.Println("Ticker stopped, deleting message")
		TeamsMessageQueue.Delete(teams_api.GUID)
	})
	logger.Printf("Starting Teams API listener on %s\n", torpedo_registry.Config.GetConfig()["teamsincomingaddr"])
	if err := http.ListenAndServe(torpedo_registry.Config.GetConfig()["teamsincomingaddr"], nil); err != nil {
		logger.Fatal(err)
	}
	tb.Stats.ConnectedAccounts -= 1
}
