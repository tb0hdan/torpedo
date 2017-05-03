package main

import (
	"fmt"
	"github.com/nlopes/slack"

	"torpedobot/multibot"
)


func HelpProcessMessage(api *slack.Client, event *slack.MessageEvent, bot *multibot.TorpedoBot) {
	idx := 0
	message := "Available commands: "
	for command := range bot.GetCommandHandlers() {
		if idx == 0 {
			message += fmt.Sprintf("!%s", command)
		} else {
			message += fmt.Sprintf(", !%s", command)
		}
		idx += 1
	}
	bot.PostMessage(event.Channel, message, api)
}
