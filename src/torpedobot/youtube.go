package main

import (
	"flag"
	"fmt"
	"strings"

	"torpedobot/multibot"
	"torpedobot/youtube"
)

var (
	google_webapp_key = flag.String("google_webapp_key", "", "Google Data API Web Application Key")
)


func YoutubeProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	message := "Usage: !youtube query\n"
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%syoutube", cmd_prefix)))
	if command != "" {
		searchResults := youtube.YoutubeSearch(command, *google_webapp_key, 25)
		message = fmt.Sprintf("https://youtu.be/%s", searchResults[0].VideoID)
	}
	bot.PostMessage(channel, message, api)
}
