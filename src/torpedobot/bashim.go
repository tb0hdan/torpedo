package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/nlopes/slack"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"torpedobot/common"
)

func windows1251_to_utf8(reader_in io.Reader) (reader_out *transform.Reader) {
	reader_out = transform.NewReader(reader_in, charmap.Windows1251.NewDecoder())
	return
}

func windows1251_to_utf8_bytes_reader(input []byte) (output *bytes.Reader, err error) {
	// TODO: Well, this is obviusly effed up
	result, err := ioutil.ReadAll(windows1251_to_utf8(bytes.NewReader(input)))
	if err != nil {
		fmt.Printf("windows1251_to_utf8_bytes_reader failed with %+v", err)
	}
	output = bytes.NewReader(result)
	return
}

func get_html(url string) (result *html.Node) {
	res, err := common.GetURLBytes(url)
	if err != nil {
		return
	}
	reader, err := windows1251_to_utf8_bytes_reader(res)

	result, err = html.Parse(reader)
	if err != nil {
		fmt.Printf("Error %s parsing html", err)
	}
	return
}

func BashProcessMessage(api *slack.Client, event *slack.MessageEvent) {
	r := get_html("http://bash.im/random")
	quotes := make(map[int]string)

	i := 0

	var f func(*html.Node)
	var quote string

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			attr := n.Attr[0]
			if attr.Key == "class" && attr.Val == "text" {
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
