package main


import (
        "bytes"
        "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "golang.org/x/net/html"
        "golang.org/x/text/transform"
        "golang.org/x/text/encoding/charmap"
        "github.com/nlopes/slack"
       )


func windows1251_to_utf8(reader_in io.ReadCloser) (reader_out *transform.Reader){
    reader_out = transform.NewReader(reader_in, charmap.Windows1251.NewDecoder())
    return
}


func get_html(url string) (result *html.Node) {
    response, err := http.DefaultClient.Get(url)
    if err != nil {
        fmt.Printf("Error getting data from %s", url)
    }
    defer response.Body.Close()

    res, _ := ioutil.ReadAll(windows1251_to_utf8(response.Body))

    result, err = html.Parse(bytes.NewReader(res))
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
