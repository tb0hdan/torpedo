package stackoverflow

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"torpedobot/common"

	"github.com/PuerkitoBio/goquery"
)

const SearchURL = "https://www.google.com/search?q=site:stackoverflow.com/questions"

type StackOverflowClient struct {
	SearchURL string
	logger    *log.Logger
}

func (so *StackOverflowClient) Search(query string) (result string, err error) {
	links := make([]string, 0)
	answers := make([]string, 0)

	cu := common.Utils{}

	search_url := fmt.Sprintf("%s%s%s", so.SearchURL, url.PathEscape(" "), url.QueryEscape(query))
	data, err := cu.GetURLBytes(search_url)
	if err != nil {
		so.logger.Printf("HTTP(S) Request to %s failed with: %+v\n", search_url, err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		so.logger.Printf("Goquery parse failed with: %+v\n", err)
		return
	}

	doc.Find("h3.r").Each(func(i int, s *goquery.Selection) {
		if attr, exists := s.Find("a").Attr("href"); exists {
			values, err := url.ParseQuery(attr)
			if err != nil {
				so.logger.Printf("doc.Find failed with: %+v\n", err)
				return
			}
			links = append(links, values["/url?q"][0])
		}
	})
	if len(links) == 0 {
		return
	}
	post_url := links[0] + "?answertab=votes#tab-top"
	data, err = cu.GetURLBytes(post_url)
	if err != nil {
		so.logger.Printf("HTTP(S) Request to %s failed with: %+v\n", post_url, err)
		return
	}
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		so.logger.Printf("Goquery parse failed with: %+v\n", err)
		return
	}
	doc.Find("td.answercell pre code").Each(func(i int, s *goquery.Selection) {
		answers = append(answers, s.Text())
	})
	if len(answers) > 0 {
		result = answers[0]
	}
	return
}

func NewClient(CustomSearchURL string) (client *StackOverflowClient) {
	cu := common.Utils{}
	client = &StackOverflowClient{}
	client.logger = cu.NewLog("stack-overflow")
	if CustomSearchURL == "" {
		client.SearchURL = SearchURL
	} else {
		client.SearchURL = CustomSearchURL
	}
	return
}
