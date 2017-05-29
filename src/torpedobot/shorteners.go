package main


import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"torpedobot/common"
	"torpedobot/multibot"
	"log"
)


func QREncoderProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	cu := &common.Utils{}
	cu.SetLogger(log.New(os.Stdout, "shorteners-plugin: ", log.Lshortfile|log.LstdFlags))
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%sqr", api.CommandPrefix)))

	if command == "" {
		api.Bot.PostMessage(channel, fmt.Sprintf("Usage: %sqr query\n", api.CommandPrefix), api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%sqr", api.CommandPrefix)))
		richmsg := multibot.RichMessage{ImageURL:fmt.Sprintf("http://chart.apis.google.com/chart?cht=qr&chs=350x350&chld=M|2&chl=%s", command), Text:command}
		api.Bot.PostMessage(channel, "", api, richmsg)
	}
}


func TinyURLProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	cu := &common.Utils{}
	cu.SetLogger(log.New(os.Stdout, "shorteners-plugin: ", log.Lshortfile|log.LstdFlags))
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%stinyurl", api.CommandPrefix)))

	if command == "" {
		api.Bot.PostMessage(channel, fmt.Sprintf("Usage: %stinyurl url\n", api.CommandPrefix), api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%stinyurl", api.CommandPrefix)))
		query := url.QueryEscape(command)
		result, err := cu.GetURLBytes(fmt.Sprintf("http://tinyurl.com/api-create.php?url=%s", query))
		message := "An error occured during TinyURL encoding process"
		if err == nil {
			message = string(result)
		}
		api.Bot.PostMessage(channel, message, api)
	}
}


func CryptoProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	requestedFeature, command, message := common.GetRequestedFeature(incoming_message)
	if command != "" {
		switch requestedFeature {
		case fmt.Sprintf("%sb64e", api.CommandPrefix):
			message = base64.StdEncoding.EncodeToString([]byte(command))
		case fmt.Sprintf("%sb64d", api.CommandPrefix):
			decoded, err := base64.StdEncoding.DecodeString(command)
			if err != nil {
				message = fmt.Sprintf("%v", err)
			} else {
				message = string(decoded)
			}
		case fmt.Sprintf("%smd5", api.CommandPrefix):
			message = common.MD5Hash(command)
		case fmt.Sprintf("%ssha1", api.CommandPrefix):
			message = common.SHA1Hash(command)
		case fmt.Sprintf("%ssha256", api.CommandPrefix):
			message = common.SHA256Hash(command)
		case fmt.Sprintf("%ssha512", api.CommandPrefix):
			message = common.SHA512Hash(command)
		default:
			// should never get here
			message = "Unknown feature requested"
		}
	}
	api.Bot.PostMessage(channel, message, api)
}
