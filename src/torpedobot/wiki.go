package main


import (
	"fmt"
	"net/url"
	"strings"

	"torpedobot/multibot"
	"torpedobot/wiki"
)


func WikiProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var richmsg multibot.RichMessage
	client := wiki.NewClient()
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%swiki", api.CommandPrefix)))
	message := fmt.Sprintf("Usage: %swiki query\n", api.CommandPrefix)
	if command != "" {
		message = "The page you've requested could not be found."
		summary := client.GetWikiPageExcerpt(command)
		if summary != "" {
			message = ""
			image_url, _ := client.GetWikiTitleImage(command)
			richmsg = multibot.RichMessage{BarColor: "#36a64f",
				 			Text: summary,
							Title: command,
							TitleLink: fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.QueryEscape(command)),
							ImageURL: image_url,}
		}
	}
	api.Bot.PostMessage(channel, message, api, richmsg)
}
