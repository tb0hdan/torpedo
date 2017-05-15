package multibot

import (
	"net/url"
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
	"strings"
	"log"
	"os"
	"time"
	"regexp"
	"torpedobot/common"
)

type SkypeIncomingMessage struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Timestamp string `json:"timestamp"`
	LocalTimestamp string `json:"localTimestamp"`
	ID string `json:"id"`
	ChannelID string `json:"channelId"`
	ServiceURL string `json:"serviceUrl"`
	From struct {
		ID string `json:"id"`
		Name string `json:"name"`
	} `json:"from"`
	Conversation struct {
		ID string `json:"id"`
	} `json:"conversation"`
	Recipient struct {
		ID string `json:"id"`
		Name string `json:"name"`
	} `json:"recipient"`
	Entities []struct {
		Locale string `json:"locale"`
		Country string `json:"country"`
		Platform string `json:"platform"`
		Type string `json:"type"`
	} `json:"entities"`
	ChannelData struct {
		Text string `json:"text"`
	} `json:"channelData"`
}

type SkypeAttachment struct {
	// base64 encoded content of media, this or ContentURL
	// data:image/png;base64,iVBORw0KGgo…
	Content string `json:"content,omitempty"`
	ContentType string `json:"contentType"`
	ContentURL string `json:"contentUrl"`
	Name string `json:"name"`
}

type SkypeOutgoingMessage struct {
	Text string `json:"text"`
	Type string `json:"type"`
	TextFormat string `json:"textFormat"`
	Attachments []*SkypeAttachment `json:"attachments,omitempty"`
}

type SkypeTokenResponse struct {
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	AccessToken string `json:"access_token"`

}


type SkypeAPI struct {
	ServiceURL string
	AccessToken string
	ExpiresIn int64
	logger *log.Logger
}


func (sapi *SkypeAPI) Send(channel, message string, attachments...*SkypeAttachment) {
	client := &http.Client{}
	outgoing_message := &SkypeOutgoingMessage{Text: message,
		Type: "message",
		TextFormat: "plain",
		Attachments: attachments,}
	parsed, _ := url.Parse(sapi.ServiceURL)
	host := parsed.Host
	body, _ := json.Marshal(outgoing_message)

	req, err := http.NewRequest("POST",
				     fmt.Sprintf("https://%s/v3/conversations/%s/activities", host, channel),
				     bytes.NewReader(body))
	sapi.logger.Printf(sapi.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sapi.AccessToken))
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


func (sapi *SkypeAPI) GetToken(app_id, app_password string) (token_response *SkypeTokenResponse){
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", app_id)
	form.Add("client_secret", app_password)
	form.Add("scope", "https://api.botframework.com/.default")

	r, err := http.DefaultClient.Post("https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token",
	"application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		sapi.logger.Printf("%+v\n", err)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	token_response = &SkypeTokenResponse{}
	err = json.Unmarshal(data, token_response)
	if err != nil {
		sapi.logger.Printf("An error occured during token unmarshalling: %+v\n", err)
	}
	return
}

func (tb *TorpedoBot) RunSkypeBot(apiKey, cmd_prefix string) {
	skype_api := &SkypeAPI{}
	logger := log.New(os.Stdout, "skype-bot: ", log.Lshortfile|log.LstdFlags)
	skype_api.logger = logger
	app_id := strings.Split(apiKey, ":")[0]
	app_password := strings.Split(apiKey, ":")[1]
	logger.Printf("Waiting for Skype token...\n")
	token_response := skype_api.GetToken(app_id, app_password)
	logger.Printf("Got Token: %s\n", token_response.AccessToken)
	skype_api.AccessToken = token_response.AccessToken
	skype_api.ExpiresIn = int64(time.Now().Unix()) + int64(token_response.ExpiresIn)

	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		defer r.Body.Close()
		body_bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			tb.logger.Fatalf("readAll errored with %+v\n", err)
			return
		}
		logger.Printf("Skype incoming message: %s\n", string(body_bytes))
		message := &SkypeIncomingMessage{}
		err = json.Unmarshal(body_bytes, message)
		if err != nil {
			logger.Fatalf("JSON unmarshalling failed with %+v\n", err)
			return
		}

		// Check token (ExpiresIn is in the future)
		if 1+skype_api.ExpiresIn-int64(time.Now().Unix()) <= 0 {
			// Get new token
			token_response := skype_api.GetToken(app_id, app_password)
			logger.Printf("Got Token: %s\n", token_response.AccessToken)
			skype_api.AccessToken = token_response.AccessToken
			skype_api.ExpiresIn = int64(time.Now().Unix()) + int64(token_response.ExpiresIn)
		} else {
			logger.Printf("Token expires in %vs\n", skype_api.ExpiresIn-int64(time.Now().Unix()))
		}

		botApi := &TorpedoBotAPI{}
		skype_api.ServiceURL = message.ServiceURL
		botApi.API = skype_api
		botApi.Bot = tb
		botApi.CommandPrefix = cmd_prefix
		re := regexp.MustCompile(`^(@(.+)\s)?`)
		msg := re.ReplaceAllString(message.Text, "")
		logger.Printf("Message: `%s`\n", msg)
		go tb.processChannelEvent(botApi, message.Conversation.ID, msg)
	})
	logger.Printf("Starting Skype API listener on %s\n", tb.config.SkypeIncomingAddr)
	http.ListenAndServe(tb.config.SkypeIncomingAddr, nil)
}
