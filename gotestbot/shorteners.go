package main

import (
       "fmt"
       "io/ioutil"
       "log"
       "os"
       "strings"

       "github.com/nlopes/slack"
       )


func QREncoderProcessMessage(api *slack.Client, event *slack.MessageEvent) {
    command := strings.Trim(strings.TrimLeft(event.Text, "!qr"), " ")

    if command == "" {
        postMessage(event.Channel, "Usage: !qr query\n", api)
    } else {
        command := strings.Trim(strings.TrimLeft(event.Text, "!qr"), " ")

        img, _ := GetURLBytes(fmt.Sprintf("http://chart.apis.google.com/chart?cht=qr&chs=350x350&chld=M|2&chl=%s", command))
        tmpfile, err := ioutil.TempFile("/tmp", "torpedo")
        if err != nil {
            log.Fatal(err)
        }

        defer os.Remove(tmpfile.Name())

        if _, err := tmpfile.Write(img); err != nil {
            log.Fatal(err)
        }

        if err := tmpfile.Close(); err != nil {
            log.Fatal(err)
        }

        channels := []string{event.Channel}
        filename := fmt.Sprintf("%s.png", command)
        parameters := slack.FileUploadParameters{File: tmpfile.Name(), Filetype: "image/png", Filename: filename, Title: filename, Channels: channels}
        api.UploadFile(parameters)
    }
}
