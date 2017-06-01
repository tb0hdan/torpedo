package main


import (
	"fmt"

	"torpedobot/games"
	"torpedobot/multibot"
)


func SteamProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	client := games.NewClient("steam", games.SteamStoreURL)
	for _, item := range client.SteamShowNew() {
		platforms := ""
		for _, platform := range item.Platforms {
			platforms += fmt.Sprintf("%s ", platform)
		}
		message += fmt.Sprintf("%s [%s] - Regular Price: %v, Current Price: %v, Difference: %v%%\n", item.GameURL, platforms, item.RegularPrice, item.CurrentPrice, item.DiscountPercentage)
	}
	api.Bot.PostMessage(channel, message, api)
}
