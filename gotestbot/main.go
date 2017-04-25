package main

import (
        "flag"
	"fmt"
	"log"
	"os"
        "strconv"
        "strings"
        "time"
	"github.com/nlopes/slack"
)

var (
    token = flag.String("token", "", "Slack legacy token")
    commandHandlers = make(map[string]func(*slack.Client, *slack.MessageEvent))
)

func postMessage(channel, message string, api *slack.Client) {
    params := slack.PostMessageParameters{}
    channelID, timestamp, err := api.PostMessage(channel, message, params)
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }
    fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func processChannelEvent(api *slack.Client, event *slack.MessageEvent) {
    messageTS, _ := strconv.ParseFloat(event.Timestamp, 64)
    jitter := int64(time.Now().Unix()) - int64(messageTS)

    if jitter < 10 && strings.HasPrefix(event.Text, "!") {
        command := strings.TrimPrefix(event.Text, "!")
        found := 0
        for handler := range commandHandlers {
            if strings.HasPrefix(command, handler) {
                found += 1
                commandHandlers[handler](api, event)
            }
        }
        fmt.Printf("PROCESS! -> %s", command)
        if found == 0 {
            postMessage(event.Channel, "Could not process your message: " + command + ". Command unknown", api)
        }
    }
}


func main() {
        flag.Parse()
	api := slack.New(*token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

        RegisterChatHandlers(commandHandlers)

        // TODO: Move this somewhere else
	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			//fmt.Printf("Message: %v\n", ev)
                        go processChannelEvent(api, ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			//fmt.Printf("Unexpected: %v\n", msg.Data)
		}
    }
}
