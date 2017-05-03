package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/nlopes/slack"

	"torpedobot/common"
	"torpedobot/file"
	"torpedobot/multibot"
)

func GetSetImageProcessMessage(api *slack.Client, event *slack.MessageEvent, bot *multibot.TorpedoBot) {
	var params slack.PostMessageParameters
	requestedFeature, command, message := common.GetRequestedFeature(event.Text)
	if command != "" {
		switch requestedFeature {
		case "!getimg":
			fpath, mimetype, err := file.GetChannelFile(event.Channel, command)
			if fpath != "" {
				common.ChannelsUploadImage([]string{event.Channel}, command, fpath, mimetype, api)
				return
			} else {
				message = fmt.Sprintf("%+v", err)
			}
		case "!setimg":
			msg, err := file.SetChannelFile(event.Channel, command)
			if err != nil {
				message = fmt.Sprintf("%+v", err)
			} else {
				message = msg
			}
		case "!listimg", "!lsimg":
			files, err := file.ListChannelFiles(event.Channel)
			if err != nil {
				message = "An error occured while retrieving image file list"
			} else {
				message = ""
				for _, fname := range files {
					msg, err := base64.URLEncoding.DecodeString(fname)
					if err != nil {
						continue
					}
					message += fmt.Sprintf("`%s`\n", msg)
				}
				if message == "" {
					message = "No files found, upload using !setimg first"
				} else {
					message = fmt.Sprintf("Available image files:\n%s", message)
				}
			}
		case "!rmimg":
			fpath, _, _ := file.GetChannelFile(event.Channel, command)
			if fpath != "" {
				err := os.Remove(fpath)
				if err != nil {
					message = fmt.Sprintf("An error occured while trying to remove `%s`", command)
				} else {
					message = fmt.Sprintf("Requested file `%s` was removed", command)
				}
			} else {
				message = fmt.Sprintf("Requested file `%s` was not found", command)
			}
		default:
			// should never get here
			message = "Unknown feature requested"
		}
	}
	bot.PostMessage(event.Channel, message, api, params)
}
