package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"torpedobot/common"
	"torpedobot/multibot"
)

func ParseFunCommand(feature, command string) (result string) {
	switch feature {
	case "rm":
		if command == "" {
			result = `rm: missing operand
Try 'rm --help' for more information.
`
		} else {
			result = "rm: removing files"
		}
	case "halt", "poweroff", "reboot", "shutdown":
		if command == "" {
			t := time.Now()
			hm := t.Format("15:04")
			result = fmt.Sprintf("Broadcast message from root@localhost\n\t(/dev/tty1) at %s ...\n\nThe system is going down for %s NOW!", hm, feature)
		} else {

		}
	case "kill":
		if command == "" {
			result = "kill: usage: kill [-s sigspec | -n signum | -sigspec] pid | jobspec ... or kill -l [sigspec]"
		} else {
			result = fmt.Sprintf("-bash: kill: (%s) - No such process", command)
		}
	default:
		result = fmt.Sprintf("-bash: %s: command not found", feature)
	}
	return
}

func FunProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	requestedFeature, command, _ := common.GetRequestedFeature(incoming_message)
	log.Printf("Feature: %s, command: %s\n", requestedFeature, command)
	message := ParseFunCommand(strings.TrimLeft(requestedFeature, api.CommandPrefix), command)
	api.Bot.PostMessage(channel, message, api)
}
