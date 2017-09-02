package multibot

import (
	"flag"
	"fmt"
	"os"

	"github.com/tb0hdan/torpedo_registry"
)

var ShowPlugins *bool

func (tb *TorpedoBot) ConfigureListPlugins(cfg *torpedo_registry.ConfigStruct) {
	ShowPlugins = flag.Bool("l", false,
		"List enabled handlers")
}

func (tb *TorpedoBot) ParseListPlugins(cfg *torpedo_registry.ConfigStruct) {
	if *ShowPlugins {
		tb.ShowCompiledInPlugins()
		os.Exit(0)
	}
}

func (tb *TorpedoBot) ShowCompiledInPlugins() {
	var (
		cnt     int
		tmp     string
		message string
	)
	// not actually plugins, rather handlers (including duplicates) but it's a start
	for name, _ := range torpedo_registry.Config.GetHandlers() {
		cnt += 1
		tmp += fmt.Sprintf("%s ", name)
	}
	if tmp != "" {
		message = fmt.Sprintf("Command handlers (%v):\n%s\n\n", cnt, tmp)
	} else {
		message = "No command handlers configured\n"
	}
	cnt = 0
	tmp = ""
	for name, _ := range torpedo_registry.Config.GetTextMessageHandlers() {
		tmp += fmt.Sprintf("%s ", name)
	}
	if tmp != "" {
		cnt += 1
		message += fmt.Sprintf("Text message handlers (%v) :\n%s\n\n", cnt, tmp)
	} else {
		message += "No text message handlers configured\n"
	}
	tb.logger.Println(message)
	return
}
