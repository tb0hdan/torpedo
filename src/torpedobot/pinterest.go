package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nlopes/slack"

	"torpedobot/common"
	"torpedobot/multibot"
	"torpedobot/pinterest"
)

var Token = flag.String("pinterest_token", "", "Pinterest Client Token")

func PinterestProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	var params slack.PostMessageParameters
	requestedFeature, command, message := common.GetRequestedFeature(incoming_message, "board")
	command = strings.Split(command, " ")[0]

	switch command {
	case "board":
		board := strings.TrimSpace(strings.TrimPrefix(incoming_message, fmt.Sprintf("%s %s", requestedFeature, command)))
		if board != "" {
			api := pinterest.New(*Token)
			images, err := api.GetImagesForBoard(board)
			if err != nil {
				return
			}
			attachment := slack.Attachment{
				Color:     "#36a64f",
				Text:      board,
				Title:     board,
				TitleLink: pinterest.PINTEREST_API_BASE + board,
				ImageURL:  images[0],
			}
			params.Attachments = []slack.Attachment{attachment}
		}
	default:
		if command != "" {
			message = fmt.Sprintf("Command %s not available yet", command)
		}
	}

	bot.PostMessage(channel, message, api, params)
}
