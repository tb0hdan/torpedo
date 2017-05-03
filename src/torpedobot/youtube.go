package main

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"

	"torpedobot/youtube"
	"flag"
)

var (
	google_webapp_key = flag.String("google_webapp_key", "", "Google Data API Web Application Key")
)


func YoutubeProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	message := "Usage: !youtube query\n"
	command := strings.TrimSpace(strings.TrimLeft(event.Text, "!youtube"))
	if command != "" {
		searchResults := youtube.YoutubeSearch(command, *google_webapp_key, 25)
		message = fmt.Sprintf("https://youtu.be/%s", searchResults[0].VideoID)
	}
	postMessage(event.Channel, message, api)
}
