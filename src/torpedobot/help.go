package main


import (
	"fmt"

	"torpedobot/multibot"
)


func HelpProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	idx := 0
	message := "Available commands: "
	for command := range api.Bot.GetCommandHandlers() {
		if idx == 0 {
			message += fmt.Sprintf("%s%s", api.CommandPrefix, command)
		} else {
			message += fmt.Sprintf(", %s%s", api.CommandPrefix, command)
		}
		idx += 1
	}
	api.Bot.PostMessage(channel, message, api)
}
