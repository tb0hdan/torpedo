package multibot

import (
	"net/url"
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	//"io/ioutil"
	"io/ioutil"
	"strings"
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


type SkypeOutgoingMessage struct {
	Text string `json:"text"`
	Type string `json:"type"`
	TextFormat string `json:"textFormat"`
}

type SkypeTokenResponse struct {
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	AccessToken string `json:"access_token"`

}


type SkypeAPI struct {
	ServiceURL string
	AccessToken string
	ExpiresIn string
}


func (sapi *SkypeAPI) Send(channel, message string) {
	client := &http.Client{}
	outgoing_message := &SkypeOutgoingMessage{Text: message,
		Type: "message",
		TextFormat: "plain"}
	parsed, _ := url.Parse(sapi.ServiceURL)
	host := parsed.Host
	body, _ := json.Marshal(outgoing_message)

	req, err := http.NewRequest("POST",
				     fmt.Sprintf("https://%s/v3/conversations/%s/activities", host, channel),
				     bytes.NewReader(body))
	fmt.Printf(sapi.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sapi.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp)
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
		fmt.Printf("%+v", err)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	token_response = &SkypeTokenResponse{}
	err = json.Unmarshal(data, token_response)
	if err != nil {
		fmt.Printf("An error occured during token unmarshalling: %+v", err)
	}
	return
}
