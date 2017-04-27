package main

import (
        "fmt"
        "flag"
        "strings"
        "github.com/nlopes/slack"
        "github.com/shkh/lastfm-go/lastfm"
       )

var (
    lastfm_key = flag.String("lastfm_key", "", "Last.FM API Key")
    lastfm_secret = flag.String("lastfm_secret", "", "Last.FM API Secret")
)

func lastfmArtist(artist string) (result string) {
    var tags string
    lastfm_api := lastfm.New(*lastfm_key, *lastfm_secret)
    r,  err := lastfm_api.Artist.GetInfo(lastfm.P{"artist": artist})
    result = "An error occured while processing your request"
    if err == nil {
        for idx, tag := range r.Tags {
            if idx == 0 {
                tags = tag.Name
            } else {
                tags += fmt.Sprintf(", %s", tag.Name)
            }
        }
        result = fmt.Sprintf("%s\nTags: %s", r.Bio.Summary, tags)
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


func LastFmProcessMessage(api *slack.Client, event *slack.MessageEvent) {
    var message string
    help := "Usage: !lastfm command\nAvailable commands: artist, tag"
    command := strings.Split(strings.Trim(strings.TrimLeft(event.Text, "!lastfm"), " "), " ")[0]

    switch command {
    case "artist":
        artist := strings.Trim(strings.TrimPrefix(event.Text, fmt.Sprintf("!lastfm %s", command)), " ")
        if artist != "" {
            message = lastfmArtist(artist)
        } else {
            message = "Please supply artist: !lastfm artist artist_name"
        }
    case "tag":
        tag := strings.Trim(strings.TrimPrefix(event.Text, fmt.Sprintf("!lastfm %s", command)), " ")
        if tag != "" {
            message = lastfmTag(tag)
        } else {
            message = "Please supply tag: !lastfm tag tag_name"
        }
    default:
        message = help
    }

    postMessage(event.Channel, message, api)
}
