package main

import (
	"fmt"

	"strings"
	"github.com/tb0hdan/torpedo_registry"
	common "github.com/tb0hdan/torpedo_common"
)

func HelpProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	var message string
	_, command, _ := common.GetRequestedFeature(incoming_message)
	if command == "" {
		idx := 0
		message = "Available commands: "
		for command := range api.Bot.GetCommandHandlers() {
			if idx == 0 {
				message += fmt.Sprintf("`%s%s`", api.CommandPrefix, command)
			} else {
				message += fmt.Sprintf(", `%s%s`", api.CommandPrefix, command)
			}
			idx += 1
		}
	} else {
		message = "No help available yet"
		for help := range api.Bot.GetHelp() {
			if strings.TrimLeft(command, api.CommandPrefix) == help {
				message = api.Bot.GetHelp()[help]
				break
			}
		}
	}
	api.Bot.PostMessage(channel, message, api)
}
