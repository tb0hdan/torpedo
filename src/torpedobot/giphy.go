package main


import (
	"torpedobot/giphy"
	"torpedobot/common"
	"torpedobot/multibot"

	"log"
	"os"
)


func GiphyProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	var richmsg multibot.RichMessage
	logger := log.New(os.Stdout, "giphy-process-message: ", log.Lshortfile|log.LstdFlags)

	client := giphy.NewClient()

	_, command, message := common.GetRequestedFeature(incoming_message)
	if command != "" {
		logger.Printf("Got command %s\n", command)
		giphyResponse := client.GiphySearch(command)
		if giphyResponse.Meta.Status == 200 {
			richmsg = multibot.RichMessage{BarColor: "#36a64f",
				                       Text: command,
			                               TitleLink: giphyResponse.Data[0].URL,
			                               ImageURL: giphyResponse.Data[0].Images.OriginalImage.URL, }
			message = ""
		} else {
			message = "Your request to Giphy could not be processed"
		}
	}
	api.Bot.PostMessage(channel, message, api, richmsg)
}
