package multibot

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
	"regexp"
	"io/ioutil"
)

var (
	TeamsIncomingAddr *string
	TeamsAPIKey        *string
)

type TeamsAPI struct {
	ServiceURL string
	logger     *log.Logger
}

func (sapi *TeamsAPI) Send(channel, message string, attachments ...*SkypeAttachment) {
	client := &http.Client{}
	outgoing_message := &SkypeOutgoingMessage{Text: message,
		Type:        "message",
		TextFormat:  "plain",
		Attachments: attachments}
	parsed, _ := url.Parse(sapi.ServiceURL)
	host := parsed.Host
	body, _ := json.Marshal(outgoing_message)

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s/v3/conversations/%s/activities", host, channel),
		bytes.NewReader(body))
	//sapi.logger.Printf(sapi.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sapi.AccessToken))
	req.Header.Set("User-Agent", common.User_Agent)
	resp, err := client.Do(req)
	if err != nil {
		sapi.logger.Printf("%+v\n", err)
		return
	}
	defer resp.Body.Close()
	sapi.logger.Println(resp)
	return
}

func HandleTeamsMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	fmt.Println("HandleTeamsMessage called...")
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
	tb.RegisteredProtocols["*multibot.TeamsAPI"] = HandleTeamsMessage
	cu := &common.Utils{}
	logger := cu.NewLog("teams-bot")

	teams_api := &TeamsAPI{}
	teams_api.logger = logger

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
		teams_api.ServiceURL = message.ServiceURL
		botApi.API = teams_api
		botApi.Bot = tb
		botApi.CommandPrefix = cmd_prefix
		botApi.UserProfile = &torpedo_registry.UserProfile{ID: message.From.ID, Nick: message.From.Name}
		// FIXME: Remove hardcode
		botApi.Me = "torpedobot"

		re := regexp.MustCompile(`^(@[^\s]+\s)?`)
		msg := re.ReplaceAllString(message.Text, "")
		logger.Printf("Message: `%s`\n", msg)
		go tb.processChannelEvent(botApi, message.Conversation.ID, msg)
	})
	logger.Printf("Starting Teams API listener on %s\n", torpedo_registry.Config.GetConfig()["teamsincomingaddr"])
	if err := http.ListenAndServe(torpedo_registry.Config.GetConfig()["teamsincomingaddr"], nil); err != nil {
		logger.Fatal(err)
	}
	tb.Stats.ConnectedAccounts -= 1
}
