package main

import (
	"fmt"

	"github.com/tb0hdan/torpedo_registry"
)

func ChatInfoProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, _ string) {
	message := fmt.Sprintf("API info:%T\n", api.ProtocolAPI)
	message += fmt.Sprintf("Chat info:%s\n", channel)
	api.Bot.PostMessage(channel, message, api)
}
