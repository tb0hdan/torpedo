package main

import (
	"fmt"

	"torpedobot/multibot"
)


func HelpProcessMessage(api *multibot.TorpedoBotAPI, bot *multibot.TorpedoBot, channel interface{}, incoming_message, cmd_prefix string) {
	idx := 0
	message := "Available commands: "
	for command := range bot.GetCommandHandlers() {
		if idx == 0 {
			message += fmt.Sprintf("%s%s", cmd_prefix, command)
		} else {
			message += fmt.Sprintf(", %s%s", cmd_prefix, command)
		}
		idx += 1
	}
	bot.PostMessage(channel, message, api)
}
