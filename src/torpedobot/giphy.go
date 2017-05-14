package main


import (
	"github.com/nlopes/slack"

	"torpedobot/giphy"
	"torpedobot/common"
	"torpedobot/multibot"

)


func GiphyProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	var params slack.PostMessageParameters

	client := giphy.NewClient()

	_, command, message := common.GetRequestedFeature(incoming_message)
	if command != "" {
		giphyResponse := client.GiphySearch(command)
		if giphyResponse.Meta.Status == 200 {
			attachment := slack.Attachment{
				Color:     "#36a64f",
				Title:     command,
				TitleLink: giphyResponse.Data[0].URL,
				ImageURL:  giphyResponse.Data[0].Images.OriginalImage.URL,
			}
			params.Attachments = []slack.Attachment{attachment}
			message = ""
		} else {
			message = "Your request to Giphy could not be processed"
		}
	}
	api.Bot.PostMessage(channel, message, api, params)
}
