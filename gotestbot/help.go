package main

import (
        "fmt"
        "github.com/nlopes/slack"
       )


func HelpProcessMessage(api *slack.Client, event *slack.MessageEvent) {
        idx := 0
        message := "Available commands: "
        for command := range commandHandlers {
            if idx == 0 {
                message += fmt.Sprintf("!%s", command)
            } else {
                message += fmt.Sprintf(", !%s", command)
            }
            idx += 1
        }
        postMessage(event.Channel, message, api)
}
