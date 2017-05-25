package main


import (
	"fmt"
	"strings"

	"torpedobot/multibot"

	"github.com/shkh/lastfm-go/lastfm"
)




type LastFmWrapper struct {
	LastFmKey string
	LastFmSecret string
}


func (lfw *LastFmWrapper) LastfmArtist(artist string) (summary, artist_url, artist_corrected, image_url string) {
	var tags string
	lastfm_api := lastfm.New(lfw.LastFmKey, lfw.LastFmSecret)
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


func (lfw *LastFmWrapper) LastfmTag(tag string) (result string) {
	var artists string
	lastfm_api := lastfm.New(lfw.LastFmKey, lfw.LastFmSecret)
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


func (lfw *LastFmWrapper) LastfmUser(user string) (result string) {
	lastfm_api := lastfm.New(lfw.LastFmKey, lfw.LastFmSecret)
	r, err := lastfm_api.User.GetInfo(lastfm.P{"user": user})
	result = "An error occured while processing your request"
	if err == nil {
		result = fmt.Sprintf("Profile information for: %s\n", r.Url)
		result += fmt.Sprintf("Play count: %+v track(s)\n", r.PlayCount)
		result += fmt.Sprintf("\nTop artists:\n")
		r2, _ := lastfm_api.User.GetTopArtists(lastfm.P{"user": user, "limit": 10})
		for idx, artist := range r2.Artists {
			result += fmt.Sprintf("%+v - %s - %s play(s)\n", idx + 1, artist.Name, artist.PlayCount)
		}
		result += fmt.Sprintf("\nTop tracks:\n")
		r3, _ := lastfm_api.User.GetTopTracks(lastfm.P{"user": user, "limit": 10})
		for idx, track := range r3.Tracks {
			result += fmt.Sprintf("%+v - %s - %s - %s play(s)\n", idx + 1, track.Artist.Name, track.Name, track.PlayCount)
		}
	}
	return
}


func LastFmProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	var richmsg multibot.RichMessage
	lfm := &LastFmWrapper{LastFmKey:api.Bot.Config.LastFmKey, LastFmSecret:api.Bot.Config.LastFmSecret}
	help := fmt.Sprintf("Usage: %slastfm command\nAvailable commands: artist, tag, user", api.CommandPrefix)
	command := strings.Split(strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%slastfm", api.CommandPrefix))), " ")[0]

	switch command {
	case "artist":
		artist := strings.TrimSpace(strings.TrimPrefix(incoming_message, fmt.Sprintf("%slastfm %s", api.CommandPrefix, command)))
		if artist != "" {
			summary, artist_url, artist_corrected, image_url := lfm.LastfmArtist(artist)
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
			message = lfm.LastfmTag(tag)
		} else {
			message = fmt.Sprintf("Please supply tag: %slastfm tag tag_name", api.CommandPrefix)
		}
	case "user":
		user := strings.TrimSpace(strings.TrimPrefix(incoming_message, fmt.Sprintf("%slastfm %s", api.CommandPrefix, command)))
		if user != "" {
			message = lfm.LastfmUser(user)
		} else {
			message = fmt.Sprintf("Please supply user name: %slastfm user user_name", api.CommandPrefix)
		}
	default:
		message = help
	}

	api.Bot.PostMessage(channel, message, api, richmsg)
}
