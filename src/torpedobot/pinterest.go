package main


import (
	"fmt"
	"strings"

	"torpedobot/common"
	"torpedobot/multibot"
	"torpedobot/pinterest"
)


func PinterestProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var richmsg multibot.RichMessage

	requestedFeature, command, message := common.GetRequestedFeature(incoming_message, "board")
	command = strings.Split(command, " ")[0]

	switch command {
	case "board":
		board := strings.TrimSpace(strings.TrimPrefix(incoming_message, fmt.Sprintf("%s %s", requestedFeature, command)))
		if board != "" {
			api := pinterest.New(api.Bot.Config.PinterestToken)
			images, err := api.GetImagesForBoard(board)
			if err != nil {
				return
			}
			richmsg = multibot.RichMessage{BarColor: "#36a64f",
			                               Text: board,
						       Title: board,
						       TitleLink: pinterest.PINTEREST_API_BASE + board,
			                               ImageURL: images[0],}
		}
	default:
		if command != "" {
			message = fmt.Sprintf("Command %s not available yet", command)
		}
	}

	api.Bot.PostMessage(channel, message, api, richmsg)
}
