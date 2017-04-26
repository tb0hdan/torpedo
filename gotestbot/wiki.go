package main

import (
        "encoding/json"
        "fmt"
        "net/url"
        "regexp"
        "strings"

        "github.com/nlopes/slack"
       )

type WikiRevision struct {
    ContentFormat string `json:"contentformat"`
    ContentModel string `json:"contentmodel"`
    LatestRevision string `json:"*"`
}


type WikiPage struct {
    PageID int `json:"pageid"`
    NS int `json:"ns"`
    Title string `json:"title"`
    Revisions []*WikiRevision `json:"revisions"`
}


type WikiQuery struct {
    Pages map[int]*WikiPage `json:"pages"`
}


type WikiResponse struct {
    BatchComplete string `json:"batchcomplete"`
    Query *WikiQuery `json:"query"`
}


func GetWikiPage(query string) (result string, err error) {
    var wikiResponse WikiResponse
    new_query := url.QueryEscape(query)
    data, _ := GetURLBytes(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=json&titles=%s", new_query))
    err = json.Unmarshal(data, &wikiResponse)
    if err != nil {
        return
    }
    for page := range wikiResponse.Query.Pages {
        result = wikiResponse.Query.Pages[page].Revisions[0].LatestRevision
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

func WikiProcessMessage(api *slack.Client, event *slack.MessageEvent) {
    var message string
    help := "Usage: !wiki query\n"
    command := strings.Trim(strings.TrimLeft(event.Text, "!wiki"), " ")
    if command == "" {
        message = help
    } else {
        message = GetWikiPageExcerpt(command)
    }
    postMessage(event.Channel, message, api)
}
