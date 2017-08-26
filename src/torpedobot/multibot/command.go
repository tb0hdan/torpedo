package multibot

import (
	"fmt"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/tb0hdan/torpedo_registry"
)

func (tb *TorpedoBot) ProcessCommandMessage(api *TorpedoBotAPI, channel interface{}, incoming_message string) {
	var chat_message string
	tb.Stats.ProcessedMessages += 1
	// is it good idea to store it here?
	// TODO: find better way
	tb.Stats.ProcessedMessagesTotal = tb.Database.GetUpdateTotalMessages(1)
	//
	command := strings.TrimPrefix(incoming_message, api.CommandPrefix)
	botapi := tb.GetBotAPI(api, channel, incoming_message)
	found := 0
	tb.logger.Printf("PROCESS! -> `%s`", command)
	handlers := torpedo_registry.Config.GetHandlers()
	for handler := range handlers {
		if strings.ToLower(strings.Split(command, " ")[0]) == handler {
			found += 1
			if torpedo_registry.Config.GetConfig()["raven"] == "yes" {
				raven.CapturePanicAndWait(func() {
					handlers[handler](botapi, channel, incoming_message)
				}, nil)
			} else {
				handlers[handler](botapi, channel, incoming_message)
			}
			break
		}
	}
	if found == 0 {
		if torpedo_registry.Config.GetConfig()["trpe_host"] != "" {
			tb.logger.Printf("Using TRPE! -> `%s`", command)
			err, result := tb.processViaTRPE(channel, incoming_message, api.CommandPrefix, torpedo_registry.Config.GetConfig()["trpe_host"])
			if err == nil {
				chat_message = result
			} else {
				chat_message = fmt.Sprintf("Could not forward message to TRPE host: %+v\n", err)
			}
		} else {
			chat_message = "Could not process your message: %s%s. Command unknown. "
			chat_message += "Send `%shelp` for list of valid commands and `%shelp command` for details."
			chat_message = fmt.Sprintf(chat_message, api.CommandPrefix, command, api.CommandPrefix, api.CommandPrefix)
		}
		api.PostMessage(channel, chat_message)
	}
	return
}
