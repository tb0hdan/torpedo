package multibot

import (
	"torpedobot/common"

	"github.com/nlopes/slack"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)


type RichMessage struct {
	BarColor  string
	Text      string
	Title     string
	TitleLink string
	ImageURL  string
}

func (rm *RichMessage) ToSlackAttachment() (params slack.PostMessageParameters) {
	attachment := slack.Attachment{
		Color:     rm.BarColor,
		Text:      rm.Text,
		Title:     rm.Title,
		TitleLink: rm.TitleLink,
		ImageURL:  rm.ImageURL,
	}
	params.Attachments = []slack.Attachment{attachment}
	return
}

func (rm *RichMessage) ToTelegramAttachment(channel int64) (msg tgbotapi.Chattable, fname string) {
	cu := &common.Utils{}
	fname, _, is_image, err := cu.DownloadToTmp(rm.ImageURL)
	if is_image && err == nil {
		msg = tgbotapi.NewPhotoUpload(channel, fname)
	}
	return
}

func (rm *RichMessage) ToSkypeAttachment() (msg SkypeOutgoingMessage) {
	return
}
