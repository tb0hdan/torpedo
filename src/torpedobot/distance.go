package main

import (
	"torpedobot/multibot"
	"torpedobot/common"
	"log"
	"googlemaps.github.io/maps"
	"golang.org/x/net/context"
	"strings"
	"fmt"
)

func DistanceProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	message := fmt.Sprintf("Usage: `%sdistance city_A city_B` or `%sdistance city_A,country_A city_B,country_B`\n", api.CommandPrefix, api.CommandPrefix)
	_, command, _ := common.GetRequestedFeature(incoming_message)
	if command != "" && len(strings.Split(command, " ")) == 2 {
		c, err := maps.NewClient(maps.WithAPIKey(api.Bot.Config.GoogleWebAppKey))
		if err != nil {
			// Okay, fatal here...
			log.Fatalf("fatal error: %+v\n", err)
		}
		r := &maps.DirectionsRequest{
			Origin:      strings.Split(command, " ")[0],
			Destination: strings.Split(command, " ")[1],
		}
		resp, _, err := c.Directions(context.Background(), r)
		if err != nil {
			log.Printf("fatal error: %+v\n", err)
			message = "Start / Destination could not be processed"
			api.Bot.PostMessage(channel, message, api)
			return
		}
		for _, item := range resp {
			message = fmt.Sprintf("Roads: %s\n", item.Summary)
			for _, lg := range item.Legs {
				message += fmt.Sprintf("Duration: %s\n", lg.Duration)
				message += fmt.Sprintf("Distance: %s\n", lg.Distance.HumanReadable)
			}
		}
	}
	api.Bot.PostMessage(channel, message, api)
}
