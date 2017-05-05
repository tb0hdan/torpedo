package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"torpedobot/common"
	"torpedobot/multibot"
)

func QREncoderProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, "!qr"))

	if command == "" {
		bot.PostMessage(channel, fmt.Sprint("Usage: %sqr query\n", cmd_prefix), api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(incoming_message, "!qr"))
		filepath, mimetype, _, _ := common.DownloadToTmp(fmt.Sprintf("http://chart.apis.google.com/chart?cht=qr&chs=350x350&chld=M|2&chl=%s", command))
		defer os.Remove(filepath)
		channels := []string{channel.(string)}
		filename := fmt.Sprintf("%s.png", command)
		common.ChannelsUploadImage(channels, filename, filepath, mimetype, api)
	}
}

func TinyURLProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%stinyurl", cmd_prefix)))

	if command == "" {
		bot.PostMessage(channel, fmt.Sprintf("Usage: %stinyurl url\n", cmd_prefix), api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%stinyurl", cmd_prefix)))
		query := url.QueryEscape(command)
		result, err := common.GetURLBytes(fmt.Sprintf("http://tinyurl.com/api-create.php?url=%s", query))
		message := "An error occured during TinyURL encoding process"
		if err == nil {
			message = string(result)
		}
		bot.PostMessage(channel, message, api)
	}
}

func CryptoProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	requestedFeature, command, message := common.GetRequestedFeature(incoming_message)
	if command != "" {
		switch requestedFeature {
		case fmt.Sprintf("%sb64e", cmd_prefix):
			message = base64.StdEncoding.EncodeToString([]byte(command))
		case fmt.Sprintf("%sb64d", cmd_prefix):
			decoded, err := base64.StdEncoding.DecodeString(command)
			if err != nil {
				message = fmt.Sprintf("%v", err)
			} else {
				message = string(decoded)
			}
		case fmt.Sprintf("%smd5", cmd_prefix):
			my_hash := md5.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		case fmt.Sprintf("%ssha1", cmd_prefix):
			my_hash := sha1.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		case fmt.Sprintf("%ssha256", cmd_prefix):
			my_hash := sha256.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		case fmt.Sprintf("%ssha512", cmd_prefix):
			my_hash := sha512.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		default:
			// should never get here
			message = "Unknown feature requested"
		}
	}
	bot.PostMessage(channel, message, api)
}
