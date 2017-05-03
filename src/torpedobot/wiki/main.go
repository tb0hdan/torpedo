package wiki

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"torpedobot/common"
)

type WikiRevision struct {
	ContentFormat  string `json:"contentformat"`
	ContentModel   string `json:"contentmodel"`
	LatestRevision string `json:"*"`
}

type WikiThumbnail struct {
	Source string `json:"source"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type WikiPage struct {
	PageID    int             `json:"pageid"`
	NS        int             `json:"ns"`
	Title     string          `json:"title"`
	Revisions []*WikiRevision `json:"revisions"`
	Thumbnail *WikiThumbnail  `json:"thumbnail,omitempty"`
}

type WikiQuery struct {
	Pages map[int]*WikiPage `json:"pages"`
}

type WikiResponse struct {
	BatchComplete string     `json:"batchcomplete"`
	Query         *WikiQuery `json:"query"`
}

func GetWikiPage(query string) (result string, err error) {
	var wikiResponse WikiResponse
	new_query := url.QueryEscape(query)
	data, _ := common.GetURLBytes(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=json&titles=%s", new_query))
	err = json.Unmarshal(data, &wikiResponse)
	if err != nil {
		return
	}

	for page := range wikiResponse.Query.Pages {
		revisions := wikiResponse.Query.Pages[page].Revisions
		if len(revisions) > 0 {
			result = revisions[0].LatestRevision
		}
	}
	return
}

func GetWikiTitleImage(query string) (result string, err error) {
	var wikiResponse WikiResponse
	new_query := url.QueryEscape(query)
	data, _ := common.GetURLBytes(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=pageimages&format=json&pithumbsize=400&titles=%s", new_query))
	err = json.Unmarshal(data, &wikiResponse)
	if err != nil {
		return
	}

	for page := range wikiResponse.Query.Pages {
		thumbnail := wikiResponse.Query.Pages[page].Thumbnail
		if thumbnail != nil {
			result = thumbnail.Source
		}
	}
	return
}

func GetWikiPageExcerpt(query string) (result string) {
	skipped := 0
	body, _ := GetWikiPage(query)
	pattern := regexp.MustCompile(`^[\[|\{|\|\*|\}|\s+\|](.+)$`)
	for _, line := range strings.Split(body, "\n") {
		if pattern.MatchString(line) {
			continue
		}
		if strings.HasPrefix(line, "==") {
			break
		}
		if line == "" {
			skipped += 1
		}
		if skipped > 1 {
			break
		}
		if line != "" {
			result += fmt.Sprintf("%s\n", line)
		}
	}
	return
}
