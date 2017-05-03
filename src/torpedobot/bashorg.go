package main

import (
	"bytes"
	"fmt"

	"github.com/nlopes/slack"
	"golang.org/x/net/html"

	"torpedobot/common"
)

func get_bashorg_html(url string) (result *html.Node) {
	res, err := common.GetURLBytes(url)
	if err != nil {
		return
	}
	result, err = html.Parse(bytes.NewReader(res))
	if err != nil {
		fmt.Printf("Error %s parsing html", err)
	}
	return
}

func BashOrgProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	r := get_bashorg_html("http://bash.org/?random")
	quotes := make(map[int]string)

	i := 0

	var f func(*html.Node)
	var quote string

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" && len(n.Attr) > 0 {
			attr := n.Attr[0]
			if attr.Key == "class" && attr.Val == "qt" {
				quote = ""
				for mc := n.FirstChild; mc != nil; mc = mc.NextSibling {
					if n.Type == html.ElementNode && mc.Data == "br" {
						continue
					}
					quote += fmt.Sprintf("%s\n", mc.Data)
				}
				quotes[i] = quote
				i += 1
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(r)

	postMessage(event.Channel, quotes[0], api)
}
