package multibot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"torpedobot/common"
)

type KikAttachment struct {
}

type KikMessage struct {
	Body       string `json:"body,omitempty"`
	PictureURL string `json:"picUrl,omitempty"`
	To         string `json:"to"`
	Type       string `json:"type"`
	ChatID     string `json:"chatId"`
}

type KikMessages struct {
	Messages []*KikMessage `json:"messages"`
}

type KikIncomingMessage struct {
	ChatID               string   `json:"chatId"`
	ID                   string   `json:"id"`
	Type                 string   `json:"type"`
	From                 string   `json:"from"`
	Participants         []string `json:"participants"`
	Body                 string   `json:"body"`
	Timestamp            int64    `json:"timestamp"`
	ReadReceiptRequested bool     `json:"readReceiptRequested"`
	Mention              string   `json:"mention"`
	PictureURL           string   `json:"picUrl,omitempty"`
	MetaData             struct {
		Product string `json:"product"`
	} `json:"metadata"`
	ChatType string `json:"chatType"`
}

type KikIncomingMessages struct {
	Messages []*KikIncomingMessage `json:"messages"`
}

type KikFeatures struct {
	ReceiveReadReceipts     bool `json:"receiveReadReceipts"`
	ReceiveIsTyping         bool `json:"receiveIsTyping"`
	ManuallySendReceipts    bool `json:"manuallySendReadReceipts"`
	ReceiveDeliveryReceipts bool `json:"receiveDeliveryReceipts"`
}

type KikAPIConfig struct {
	Webhook  string      `json:"webhook"`
	Features KikFeatures `json:"features"`
}

type KikAPI struct {
	AccessToken string
	WebHook     string
	logger      *log.Logger
}

func (ka *KikAPI) GetToken(app_id, app_password string) {
	auth_pair := fmt.Sprintf("%s:%s", app_id, app_password)
	ka.AccessToken = base64.StdEncoding.EncodeToString([]byte(auth_pair))
	return

}

func (ka *KikAPI) Configure() {
	client := &http.Client{}
	config := &KikAPIConfig{Webhook: ka.WebHook,
		Features: KikFeatures{ReceiveReadReceipts: false,
			ReceiveIsTyping:         false,
			ManuallySendReceipts:    false,
			ReceiveDeliveryReceipts: false}}
	config_json, err := json.Marshal(&config)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://api.kik.com/v1/config",
		bytes.NewReader(config_json))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", ka.AccessToken))
	req.Header.Set("User-Agent", common.User_Agent)
	resp, err := client.Do(req)
	if err != nil {
		ka.logger.Printf("%+v\n", err)
		return
	}
	defer resp.Body.Close()
	return
}

func (ka *KikAPI) SendMessages(messages *KikMessages) {
	client := &http.Client{}
	config_json, err := json.Marshal(messages)
	if err != nil {
		return
	}
	ka.logger.Printf("%s", string(config_json))
	req, err := http.NewRequest("POST", "https://api.kik.com/v1/message",
		bytes.NewReader(config_json))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", ka.AccessToken))
	req.Header.Set("User-Agent", common.User_Agent)
	resp, err := client.Do(req)
	if err != nil {
		ka.logger.Printf("%+v\n", err)
		return
	}
	defer resp.Body.Close()
	return
}

func (ka *KikAPI) Text(channel, to, message string) {
	msgs := make([]*KikMessage, 1)
	msgs[0] = &KikMessage{Body: message, To: to, Type: "text", ChatID: channel}
	messages := &KikMessages{Messages: msgs}
	ka.SendMessages(messages)
}

func (ka *KikAPI) Image(channel, to, url string) {
	msgs := make([]*KikMessage, 1)
	msgs[0] = &KikMessage{PictureURL: url, To: to, Type: "picture", ChatID: channel}
	messages := &KikMessages{Messages: msgs}
	ka.SendMessages(messages)
}

func HandleKikMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []RichMessage) {
	switch api := tba.API.(type) {
	case *KikAPI:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			msg, url := richmsgs[0].ToGenericAttachment()
			api.Text(channel.(string), tba.From, msg)
			api.Image(channel.(string), tba.From, url)
		} else {
			api.Text(channel.(string), tba.From, message)
		}
	}
}

func (tb *TorpedoBot) RunKikBot(apiKey, cmd_prefix string) {
	logger := log.New(os.Stdout, "kik-bot: ", log.Lshortfile|log.LstdFlags)
	api := &KikAPI{}
	api.logger = logger
	api.WebHook = tb.Config.KikWebHook
	api.GetToken(strings.Split(apiKey, ":")[0], strings.Split(apiKey, ":")[1])
	api.Configure()

	tb.RegisteredProtocols["*KikAPI"] = HandleKikMessage

	http.HandleFunc("/incoming", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		defer r.Body.Close()
		body_bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Fatalf("readAll failed with %+v\n", err)
			return
		}
		logger.Printf("Kik incoming message: %s\n", string(body_bytes))
		messages := &KikIncomingMessages{}
		err = json.Unmarshal(body_bytes, messages)
		if err != nil {
			logger.Fatalf("JSON unmarshalling failed with %+v\n", err)
			return
		}
		for _, message := range messages.Messages {
			botApi := &TorpedoBotAPI{}
			botApi.API = api
			botApi.Bot = tb
			botApi.CommandPrefix = cmd_prefix
			botApi.From = message.From
			logger.Printf("Message: `%s`\n", message.Body)
			go tb.processChannelEvent(botApi, message.ChatID, message.Body)
		}
	})
	logger.Printf("Starting Kik API listener on %s\n", tb.Config.KikIncomingAddr)
	if err := http.ListenAndServe(tb.Config.KikIncomingAddr, nil); err != nil {
		logger.Fatal(err)
	}
}
