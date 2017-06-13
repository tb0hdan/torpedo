package multibot

type RichMessage struct {
	BarColor  string
	Text      string
	Title     string
	TitleLink string
	ImageURL  string
}

func (rm *RichMessage) IsEmpty() bool {
	return rm.Text == "" || rm.ImageURL == ""
}

func (rm *RichMessage) ToGenericAttachment() (msg, url string) {
	msg = rm.Text
	url = rm.ImageURL
	return
}
