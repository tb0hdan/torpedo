package torpedo_text_plugin

import (
	"strings"

	"github.com/tb0hdan/torpedo_registry"
)

func TextProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	if strings.HasPrefix("hola", incoming_message) {
		api.Bot.PostMessage(channel, "buenos dias", api)
	}
}

func init() {
	torpedo_registry.Config.RegisterTextMessageHandler("text", TextProcessMessage)
}
