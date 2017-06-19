package main


import (
	"fmt"
	"strings"

	"torpedobot/multibot"
	"torpedobot/youtube"
)


func YoutubeProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	message := fmt.Sprintf("Usage: %syoutube query\n", api.CommandPrefix)
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%syoutube", api.CommandPrefix)))
	if command != "" {
		searchResults := youtube.YoutubeSearch(command, *api.Bot.Config.GoogleWebAppKey, 25)
		message = fmt.Sprintf("https://youtu.be/%s", searchResults[0].VideoID)
	}
	api.Bot.PostMessage(channel, message, api)
}
