package main

import (
	"torpedobot/multibot"
	"fmt"
	"torpedobot/stackoverflow"
)

func StrackOverflowProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	client := stackoverflow.NewClient("")
	result, err := client.Search(incoming_message)
	if err != nil {
		message = fmt.Sprintf("An error occured during StackOverflow search: %+v\n", err)
	} else {
		message = result
	}
	api.Bot.PostMessage(channel, message, api)
}

