package giphy

import (
	"log"
	"torpedobot/common"
	"net/url"
	"fmt"
	"encoding/json"
)

type GiphyUser struct {
	AvatarURL   string `json:"avatar_url"`
	BannerURL   string `json:"banner_url"`
	ProfileURL  string `json:"profile_url"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Twitter     string `json:"twitter"`
}


type GiphyFixedHeightImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	MP4      string `json:"mp4"`
	MP4Size  string `json:"mp4_size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
}


type GiphyFixedHeightStillImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyFixedHeightDownsampledImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
}


type GiphyFixedWidthImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	MP4      string `json:"mp4"`
	MP4Size  string `json:"mp4_size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
}


type GiphyFixedWidthStillImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyFixedWidthDownsampledImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
}


type GiphyFixedHeightSmallImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	MP4      string `json:"mp4"`
	MP4Size  string `json:"mp4_size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
}


type GiphyFixedHeightSmallStillImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyFixedWidthSmallImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	MP4      string `json:"mp4"`
	MP4Size  string `json:"mp4_size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
}


type GiphyFixedWidthSmallStillImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyDownsizedImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyDownsizedStillImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyDownsizedLargeImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyDownsizedMediumImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyOriginalImage struct {
	URL      string `json:"url"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	Frames   string `json:"frames"`
	MP4      string `json:"mp4"`
	MP4Size  string `json:"mp4_size"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
	Hash     string `json:"hash"`
}


type GiphyOriginalStillImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyLoopingImage struct {
	MP4     string `json:"mp4"`
	MP4Size string `json:"mp4_size"`
}


type GiphyOriginalMP4Image struct {
	MP4     string `json:"mp4"`
	MP4Size string `json:"mp4_size"`
	Width   string `json:"width"`
	Height  string `json:"height"`
}


type GiphyPreviewImage struct {
	MP4     string `json:"mp4"`
	MP4Size string `json:"mp4_size"`
	Width   string `json:"width"`
	Height  string `json:"height"`
}


type GiphyDownsizedSmallImage struct {
	MP4     string `json:"mp4"`
	MP4Size string `json:"mp4_size"`
	Width   string `json:"width"`
	Height  string `json:"height"`
}


type GiphyPreviewGIFImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyPreviewWebpImage struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}


type GiphyImages struct {
	FixedHeightImage            *GiphyFixedHeightImage            `json:"fixed_height"`
	FixedHeightStillImage       *GiphyFixedHeightStillImage       `json:"fixed_height_still"`
	FixedHeightDownsampledImage *GiphyFixedHeightDownsampledImage `json:"fixed_height_downsampled"`
	FixedWidthImage             *GiphyFixedWidthImage             `json:"fixed_width"`
	FixedWidthStillImage        *GiphyFixedWidthStillImage        `json:"fixed_width_still"`
	FixedWidthDownsampledImage  *GiphyFixedWidthDownsampledImage  `json:"fixed_width_downsampled"`
	FixedHeightSmallImage       *GiphyFixedHeightSmallImage       `json:"fixed_height_small"`
	FixedHeightSmallStillImage  *GiphyFixedHeightSmallStillImage  `json:"fixed_height_small_still"`
	FixedWidthSmallImage        *GiphyFixedWidthSmallImage        `json:"fixed_width_small"`
	FixedWidthSmallStillImage   *GiphyFixedWidthSmallStillImage   `json:"fixed_width_small_still"`
	DownsizedImage              *GiphyDownsizedImage              `json:"downsized"`
	DownsizedStillImage         *GiphyDownsizedStillImage         `json:"downsized_still"`
	DownsizedLargeImage         *GiphyDownsizedLargeImage         `json:"downsized_large"`
	DownsizedMediumImage        *GiphyDownsizedMediumImage        `json:"downsized_medium"`
	OriginalImage               *GiphyOriginalImage               `json:"original"`
	OriginalStillImage          *GiphyOriginalStillImage          `json:"original_still"`
	LoopingImage                *GiphyLoopingImage                `json:"looping"`
	OriginalMP4Image            *GiphyOriginalMP4Image            `json:"original_mp4"`
	PreviewImage                *GiphyPreviewImage                `json:"preview"`
	DownsizedSmallImage         *GiphyDownsizedSmallImage         `json:"downsized_small"`
	PreviewGIFImage             *GiphyPreviewGIFImage             `json:"preview_gif"`
	PreviewWebpImage            *GiphyPreviewWebpImage            `json:"preview_webp"`
}


type GiphyResponseDataItem struct {
	Type             string       `json:"type"`
	ID               string       `json:"id"`
	Slug             string       `json:"slug"`
	URL              string       `json:"url"`
	BitlyGIFURL      string       `json:"bitly_gif_url"`
	BitlyURL         string       `json:"bitly_url"`
	EmbedURL         string       `json:"embed_url"`
	Username         string       `json:"username"`
	Source           string       `json:"source"`
	Rating           string       `json:"rating"`
	ContentURL       string       `json:"content_url"`
	User             *GiphyUser   `json:"user"`
	SourceTLD        string       `json:"source_tld"`
	SourcePostURL    string       `json:"source_post_url"`
	IsIndexable      int          `json:"is_indexable"`
	ImportDateTime   string       `json:"import_datetime"`
	TrendingDateTime string       `json:"trending_datetime"`
	Images           *GiphyImages `json:"images"`
}


type GiphyPagination struct {
	TotalCount int `json:"total_count"`
	Count      int `json:"count"`
	Offset     int `json:"offset"`
}


type GiphyMeta struct {
	Status     int    `json:"status"`
	Message    string `json:"msg"`
	ResponseID string `json:"response_id"`
}


type GiphyResponse struct {
	Data       []*GiphyResponseDataItem `json:"data"`
	Pagination *GiphyPagination         `json:"pagination"`
	Meta       *GiphyMeta               `json:"meta"`
}


type GiphyClient struct {
	logger *log.Logger
	utils *common.Utils
}


func (gic *GiphyClient) GiphySearch(command string) (response *GiphyResponse){

	query := url.QueryEscape(command)
	result, err := gic.utils.GetURLBytes(fmt.Sprintf("https://api.giphy.com/v1/gifs/search?q=%s&api_key=dc6zaTOxFJmzC", query))
	if err != nil {
		gic.logger.Printf("Get Giphy URL failed with %+v", err)
		return
	}
	err = json.Unmarshal(result, &response)
	if err != nil {
		gic.logger.Printf("Error unmarshalling Giphy: %+v", err)
		return
	}
	return
}


func  NewClient() (client *GiphyClient){
	client = &GiphyClient{}
	client.utils = &common.Utils{}
	client.utils.SetLoggerPrefix("giphy-plugin")
	return
}
