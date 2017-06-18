package youtube


import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/google/google-api-go-client/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"torpedobot/common"
)


type YoutubeVideo struct {
	VideoID    string
	VideoTitle string
	VideoScore int
}


func (y YoutubeVideo) String() string {
	return fmt.Sprintf("%s: %d", y.VideoID, y.VideoScore)
}


type ByScore []YoutubeVideo


func (s ByScore) Len() int           { return len(s) }
func (s ByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByScore) Less(i, j int) bool { return s[i].VideoScore > s[j].VideoScore }


func CalculateVideoScore(title, query string) (score int) {
	// Rulesets
	if strings.HasSuffix(strings.ToLower(title), "(official video)") {
		score += 20
	}
	if strings.HasPrefix(strings.ToLower(title), strings.ToLower(query)) {
		score += 50
	}
	if strings.Contains(strings.ToLower(title), "(Acoustic)") {
		score -= 20
	}
	//
	return
}


func YoutubeSearch(query, developerKey string, maxResults int64) (videos []YoutubeVideo) {

	cu := &common.Utils{}
	logger := cu.NewLog("youtube-search")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		logger.Fatalf("Error creating new YouTube client: %v", err)
	}

	call := service.Search.List("id,snippet").
		Q(query).
		MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			score := CalculateVideoScore(item.Snippet.Title, query)
			videos = append(videos, YoutubeVideo{VideoID: item.Id.VideoId, VideoTitle: item.Snippet.Title, VideoScore: score})
		}
	}
	sort.Sort(ByScore(videos))
	return
}
