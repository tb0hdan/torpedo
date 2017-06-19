package main

import (
	"torpedobot/multibot"
	"torpedobot/common"
	"net/url"
)

type SoundCloudSearchResponseUser struct {
	ID int64 `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Kind string `json:"kind"`
	LastModified string `json:"last_modified"`
	Permalink string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
	URI string `json:"uri"`
	UserName string `json:"username"`
}

type SoundCloudSearchResponse struct {
	User *SoundCloudSearchResponseUser `json:"user"`
	ID int64 `json:"id"`
	AttachmentsURI string `json:"attachments_uri"`
	BPM string `json:"bpm"`
	Kind string `json:"kind"`
	CreatedAt string `json:"created_at"`
	LastModified string `json:"last_modified"`
	Permalink string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
	Title string `json:"title"`
	Duration int64 `json:"duration"`
	Sharing string `json:"sharing"`
	WaveformURL string `json:"waveform_url"`
	StreamURL string `json:"stream_url"`
	URI string `json:"uri"`
	UserID int64 `json:"user_id"`
	ArtworkURL string `json:"artwork_url"`
	CommentCount int64 `json:"comment_count"`
	Commentable bool `json:"commentable"`
	Description string `json:"description"`
	DownloadCount int64 `json:"download_count"`
	Downloadable bool `json:"downloadable"`
	EmbeddableBy string `json:"embeddable_by"`
	Favoritings_count int64 `json:"favoritings_count"`
	Genre string `json:"genre"`
	ISRC string `json:"isrc"`
	LabelID string `json:"label_id"`
	LabelName string `json:"label_name"`
	License string `json:"license"`
	LikesCount int64 `json:"likes_count"`
	OriginalContentSize int64 `json:"original_content_size"`
	OriginalFormat string `json:"original_format"`
	PlaybackCount int64 `json:"playback_count"`
	PurchaseTitle string `json:"purchase_title"`
	PurchaseURL string `json:"purchase_url"`
	Release string `json:"release"`
	ReleaseDay string `json:"release_day"`
	ReleaseMonth string `json:"release_month"`
	ReleaseYear string `json:"release_year"`
	RepostsCount int64 `json:"reposts_count"`
	State string `json:"state"`
	Streamable bool `json:"streamable"`
	TagList string `json:"tag_list"`
	TrackType string `json:"track_type"`
	KeySignature string `json:"key_signature"`
	UserFaforite bool `json:"user_favorite"`
	UserPlaybackCount int64 `json:"user_playback_count"`
	VideoURL string `json:"video_url"`
	DownloadURL string `json:"download_url"`
}


func SoundCloudProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	var message string
	cu := &common.Utils{}
	logger := cu.NewLog("soundcloud-process-message")
	_, command, _ := common.GetRequestedFeature(incoming_message)
	result := make([]SoundCloudSearchResponse, 0)
	query := url.Values{}
	query.Add("client_id", *api.Bot.Config.SoundCloudClientID)
	query.Add("q", url.QueryEscape(command))
	url := "https://api.soundcloud.com/tracks/?" + query.Encode()
	err := cu.GetURLUnmarshal(url, &result)
	if err != nil {
		message = "An error occured while searching for track on SoundCloud"
		logger.Printf("SoundCloud error: %+v, %s\n", err, url)
		api.Bot.PostMessage(channel, message, api)
		return
	}
	message = result[0].PermalinkURL
	api.Bot.PostMessage(channel, message, api)
}
