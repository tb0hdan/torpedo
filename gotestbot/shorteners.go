package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"net/url"
)


func QREncoderProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	command := strings.TrimSpace(strings.TrimLeft(event.Text, "!qr"))

	if command == "" {
		postMessage(event.Channel, "Usage: !qr query\n", api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(event.Text, "!qr"))
		filepath, mimetype, _, _ := DownloadToTmp(fmt.Sprintf("http://chart.apis.google.com/chart?cht=qr&chs=350x350&chld=M|2&chl=%s", command))
		defer os.Remove(filepath)
		channels := []string{event.Channel}
		filename := fmt.Sprintf("%s.png", command)
		ChannelsUploadImage(channels, filename, filepath, mimetype, api)
	}
}

func TinyURLProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	command := strings.TrimSpace(strings.TrimLeft(event.Text, "!tinyurl"))

	if command == "" {
		postMessage(event.Channel, "Usage: !tinyurl url\n", api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(event.Text, "!tinyurl"))
		query := url.QueryEscape(command)
		result, err := GetURLBytes(fmt.Sprintf("http://tinyurl.com/api-create.php?url=%s", query))
		message := "An error occured during TinyURL encoding process"
		if err == nil {
			message = string(result)
		}
		postMessage(event.Channel, message, api)
	}
}


func CryptoProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	requestedFeature, command, message := GetRequestedFeature(event.Text)
	if command != "" {
		switch requestedFeature {
		case "!b64e":
			message = base64.StdEncoding.EncodeToString([]byte(command))
		case "!b64d":
			decoded, err := base64.StdEncoding.DecodeString(command)
			if err != nil {
				message = fmt.Sprintf("%v", err)
			} else {
				message = string(decoded)
			}
		case "!md5":
			my_hash := md5.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		case "!sha1":
			my_hash := sha1.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		case "!sha256":
			my_hash := sha256.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		case "!sha512":
			my_hash := sha512.New()
			io.WriteString(my_hash, command)
			message = fmt.Sprintf("%x", my_hash.Sum(nil))
		default:
			// should never get here
			message = "Unknown feature requested"
		}
	}
	postMessage(event.Channel, message, api)
}
