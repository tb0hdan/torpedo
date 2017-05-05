package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nlopes/slack"

	"torpedobot/multibot"
	"torpedobot/wiki"
)

func WikiProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	var params slack.PostMessageParameters
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%swiki", cmd_prefix)))
	message := "Usage: !wiki query\n"
	if command != "" {
		message = "The page you've requested could not be found."
		summary := wiki.GetWikiPageExcerpt(command)
		if summary != "" {
			message = ""
			image_url, _ := wiki.GetWikiTitleImage(command)
			attachment := slack.Attachment{
				Color:     "#36a64f",
				Text:      summary,
				Title:     command,
				TitleLink: fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.QueryEscape(command)),
				ImageURL:  image_url,
			}
			params.Attachments = []slack.Attachment{attachment}
		}
	}
	bot.PostMessage(channel, message, api, params)
}
