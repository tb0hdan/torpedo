package main

import (
	"fmt"

	"torpedobot/steam"
	"torpedobot/multibot"
)


func SteamProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	client := steam.NewClient()
	for _, item := range client.ShowNew() {
		platforms := ""
		for _, platform := range item.Platforms {
			platforms += fmt.Sprintf("%s ", platform)
		}
		message += fmt.Sprintf("%s [%s] - Regular Price: %v, Current Price: %v, Difference: %v%%\n", item.GameURL, platforms, item.RegularPrice, item.CurrentPrice, item.DiscountPercentage)
	}
	api.Bot.PostMessage(channel, message, api)
}
