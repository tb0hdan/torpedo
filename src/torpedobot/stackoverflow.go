package main

import (
	"torpedobot/multibot"
	"fmt"
	"torpedobot/stackoverflow"
	"torpedobot/common"
	"log"
)

func StackOverflowProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	_, command, message := common.GetRequestedFeature(incoming_message)

	if command != "" {
		log.Printf("Got command %s\n", command)
		client := stackoverflow.NewClient("")
		result, err := client.Search(incoming_message)
		if err != nil {
			message = fmt.Sprintf("An error occured during StackOverflow search: %+v\n", err)
		} else {
			message = result
		}
		if message == "" {
			message = "No results for your query"
		}
	}
	api.Bot.PostMessage(channel, message, api)
}

