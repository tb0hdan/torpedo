package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nlopes/slack"

	"torpedobot/wiki"
)

func WikiProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	var params slack.PostMessageParameters
	command := strings.TrimSpace(strings.TrimLeft(event.Text, "!wiki"))
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
	postMessage(event.Channel, message, api, params)
}
