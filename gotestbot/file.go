package main

import (
        "encoding/base64"
        "fmt"
        "os"
        "strings"

        "github.com/nlopes/slack"
       )


func GetCreateChannelDir(channel string) (channelDir string, err error) {
    wd, err := os.Getwd()
    if err != nil {
        return
    }
    channelDirPath := fmt.Sprintf("%s%s%s%s%s", wd, string(os.PathSeparator), "data", string(os.PathSeparator),  channel)
    err = os.MkdirAll(channelDirPath, 0755)
    if err == nil {
        channelDir = channelDirPath
    }
    return
}


func GetChannelFile(channel, message string) (channelFile string, err error) {
    wd, err := GetCreateChannelDir(channel)
    if err != nil {
        return
    }
    // TODO: Add message permutations
    encoded := base64.URLEncoding.EncodeToString([]byte(strings.TrimSpace(message)))
    fname := fmt.Sprintf("%s%s%s", wd, string(os.PathSeparator), encoded)
    // TODO: Find a way around this, os.IsExist expects an error and we don't have one yet
    exists := true
    if _, err := os.Stat(fname); os.IsNotExist(err) {
        exists = false
    }
    if exists {
        channelFile = fname
    }
    return
}


func SetChannelFile(channel, message string) (result string, err error) {
    wd, err := GetCreateChannelDir(channel)
    if err != nil {
        return
    }
    url_formatted := strings.Split(message, " ")[0]
    url := UnformatURL(url_formatted)
    if ! (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
        result = "No valid URL found"
        return
    }
    destination := strings.TrimSpace(strings.TrimLeft(message, url_formatted))
    if destination == "" {
        result = "No valid destination found"
        return
    }
    // Check if target already exists before downloading new one
    encoded := base64.URLEncoding.EncodeToString([]byte(destination))
    new_name := fmt.Sprintf("%s%s%s", wd, string(os.PathSeparator), encoded)
    _, err = os.Open(new_name)
    if os.IsExist(err) {
        result = "Destination already exists, set skipped. Use `!rmimg destination` to remove."
        return
    }
    fname, _, is_image, err := DownloadToTmp(url)
    if is_image {
        err = os.Rename(fname, new_name)
        if err != nil {
            result = "There was an issue with setting image"
        } else {
            result = "Image set"
        }
    }
    return
}


func GetSetImageProcessMessage(api *slack.Client, event *slack.MessageEvent) {
    var params slack.PostMessageParameters
    requestedFeature, command, message := GetRequestedFeature(event.Text)
    if command != "" {
        switch requestedFeature {
        case "!getimg":
            fpath, err := GetChannelFile(event.Channel, command)
            if fpath != "" {
                fname := command
                ftype := "image/gif"
                ChannelsUploadImage([]string{event.Channel}, fname, fpath, ftype, api)
                return
            } else {
                message = fmt.Sprintf("%+v", err)
            }
        case "!setimg":
            msg, err := SetChannelFile(event.Channel, command)
            if err != nil {
                message = fmt.Sprintf("%+v", err)
            } else {
                message = msg
            }
        case "!listimg":
            message = "No action set yet"
        case "!rmimg":
            message = "No action set yet"
        default:
            // should never get here
            message = "Unknown feature requested"
        }
    }
    postMessage(event.Channel, message, api, params)
}
