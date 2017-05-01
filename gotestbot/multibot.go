package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
)


func postMessage(channel, message string, api *slack.Client, parameters...slack.PostMessageParameters) {
	var params slack.PostMessageParameters

	if len(parameters) > 0 {
		params = parameters[0]
	}
	channelID, timestamp, err := api.PostMessage(channel, message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}


func processChannelEvent(api *slack.Client, event *slack.MessageEvent, commandHandlers map[string]func(*slack.Client, *slack.MessageEvent)) {
	messageTS, _ := strconv.ParseFloat(event.Timestamp, 64)
	jitter := int64(time.Now().Unix()) - int64(messageTS)

	if jitter < 10 && strings.HasPrefix(event.Text, "!") {
		command := strings.TrimPrefix(event.Text, "!")
		found := 0
		for handler := range commandHandlers {
			if strings.HasPrefix(strings.Split(command, " ")[0], handler) {
				found += 1
				commandHandlers[handler](api, event)
				break
			}
		}
		fmt.Printf("PROCESS! -> %s", command)
		if found == 0 {
			postMessage(event.Channel, fmt.Sprintf("Could not process your message: !%s. Command unknown. Send !help for list of valid commands.", command), api)
		}
	}
}



func RunBot(apiKey string) {
	api := slack.New(apiKey)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	commandHandlers := RegisterChatHandlers()

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
			fmt.Printf("Message: %v\n", ev)
			go processChannelEvent(api, ev, commandHandlers)

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

func RunLoop() {
	for {
		time.Sleep(time.Second)
	}
}

func RunBots(keys []string) {
	for _, key := range keys {
		go RunBot(key)
	}
	RunLoop()
}
