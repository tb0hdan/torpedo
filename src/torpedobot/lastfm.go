package main


import (
	"flag"
	"fmt"
	"strings"

	"torpedobot/multibot"

	"github.com/tb0hdan/lastfm-go/lastfm"
)


var (
	lastfm_key    = flag.String("lastfm_key", "", "Last.FM API Key")
	lastfm_secret = flag.String("lastfm_secret", "", "Last.FM API Secret")
)


func lastfmArtist(artist string) (summary, artist_url, artist_corrected, image_url string) {
	var tags string
	lastfm_api := lastfm.New(*lastfm_key, *lastfm_secret)
	r, err := lastfm_api.Artist.GetInfo(lastfm.P{"artist": artist})
	summary = "An error occured while processing your request"
	if err == nil {
		for idx, tag := range r.Tags {
			if idx == 0 {
				tags = tag.Name
			} else {
				tags += fmt.Sprintf(", %s", tag.Name)
			}
		}

		for _, img := range r.Images {
			if img.Size == "large" {
				image_url = img.Url
				break
			}
		}

		summary = fmt.Sprintf("%s\n\nTags: %s\n", r.Bio.Summary, tags)
		artist_url = r.Url
		r, err := lastfm_api.Artist.GetCorrection(lastfm.P{"artist": artist})
		artist_corrected = artist
		if err == nil {
			artist_corrected = r.Correction.Artist.Name
		}
	}
	return
}

func lastfmTag(tag string) (result string) {
	var artists string
	lastfm_api := lastfm.New(*lastfm_key, *lastfm_secret)
	r, err := lastfm_api.Tag.GetTopArtists(lastfm.P{"tag": tag})
	result = "An error occured while processing your request"
	if err == nil {
		for idx, artist := range r.Artists {
			if idx == 0 {
				artists = artist.Name
			} else {
				artists += fmt.Sprintf(", %s", artist.Name)
			}
		}
		if artists != "" {
			result = fmt.Sprintf("Artists: %s\n", artists)
		}
	}
	return
}

func LastFmProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	var richmsg multibot.RichMessage
	help := fmt.Sprintf("Usage: %slastfm command\nAvailable commands: artist, tag", api.CommandPrefix)
	command := strings.Split(strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%slastfm", api.CommandPrefix))), " ")[0]

	switch command {
	case "artist":
		artist := strings.TrimSpace(strings.TrimPrefix(incoming_message, fmt.Sprintf("%slastfm %s", api.CommandPrefix, command)))
		if artist != "" {
			summary, artist_url, artist_corrected, image_url := lastfmArtist(artist)
			richmsg = multibot.RichMessage{BarColor: "#36a64f",
				Text:      summary,
				Title:     artist_corrected,
				TitleLink: artist_url,
				ImageURL:  image_url}
		} else {
			message = fmt.Sprintf("Please supply artist: %slastfm artist artist_name", api.CommandPrefix)
		}
	case "tag":
		tag := strings.TrimSpace(strings.TrimPrefix(incoming_message, fmt.Sprintf("%slastfm %s", api.CommandPrefix, command)))
		if tag != "" {
			message = lastfmTag(tag)
		} else {
			message = fmt.Sprintf("Please supply tag: %slastfm tag tag_name", api.CommandPrefix)
		}
	default:
		message = help
	}

	api.Bot.PostMessage(channel, message, api, richmsg)
}
