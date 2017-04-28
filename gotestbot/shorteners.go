package main

import (
       "crypto/md5"
       "crypto/sha1"
       "crypto/sha256"
       "crypto/sha512"
       "encoding/base64"
       "fmt"
       "io"
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


func CryptoProcessMessage(api *slack.Client, event *slack.MessageEvent) {
    requestedFeature := strings.Split(event.Text, " ")[0]
    command := strings.Trim(strings.TrimLeft(event.Text, requestedFeature), " ")
    message := fmt.Sprintf("Usage: %s string\n", requestedFeature)
    if command != "" {
        switch requestedFeature {
        case "!b64e":
            message = base64.StdEncoding.EncodeToString([]byte(command))
        case "!b64d":
            decoded, err := base64.StdEncoding.DecodeString(command)
            if err != nil {
                message = fmt.Sprintf("%v", err)
            } else {
                message = string(decoded)
            }
        case "!md5":
            my_hash := md5.New()
            io.WriteString(my_hash, command)
            message = fmt.Sprintf("%x", my_hash.Sum(nil))
        case "!sha1":
            my_hash := sha1.New()
            io.WriteString(my_hash, command)
            message = fmt.Sprintf("%x", my_hash.Sum(nil))
        case "!sha256":
            my_hash := sha256.New()
            io.WriteString(my_hash, command)
            message = fmt.Sprintf("%x", my_hash.Sum(nil))
        case "!sha512":
            my_hash := sha512.New()
            io.WriteString(my_hash, command)
            message = fmt.Sprintf("%x", my_hash.Sum(nil))
        default:
            // should never get here
            message = "Unknown feature requested"
        }
    }
    postMessage(event.Channel, message, api)
}
