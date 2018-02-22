package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/tb0hdan/torpedo_registry"
)

func GetMemStats() (result string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	result = fmt.Sprintf("\n\tAlloc = %s\n\tSys = %s\n\tNumGC = %s\n\n", humanize.Bytes(m.Alloc), humanize.Bytes(m.Sys), humanize.Comma(int64(m.NumGC)))
	return
}

func StatsProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	i, err := strconv.ParseInt(fmt.Sprintf("%v", api.Bot.Stats.StartTimestamp), 10, 64)
	if err != nil {
		panic(err)
	}
	message := fmt.Sprintf("Project URL: %s\n", api.Bot.Build.ProjectURL)
	message += fmt.Sprintf("Build hash: %s\n", api.Bot.Build.Build)
	message += fmt.Sprintf("Build date: %s\n", api.Bot.Build.BuildDate)
	message += fmt.Sprintf("Build version: %s\n", api.Bot.Build.Version)
	message += fmt.Sprintf("Golang version: %s\n", api.Bot.Build.GoVersion)
	message += fmt.Sprintf("Bot started: %s\n", humanize.Time(time.Unix(i, 0)))
	message += fmt.Sprintf("Processed messages (session): %s\n", humanize.Comma(api.Bot.Stats.ProcessedMessages))
	message += fmt.Sprintf("Processed messages (total): %s\n", humanize.Comma(api.Bot.Stats.ProcessedMessagesTotal))
	//message += fmt.Sprintf("Message rate: %s\n", CalculateMessageRate(ts, api.Bot.Stats.ProcessedMessages))
	message += fmt.Sprintf("Accounts connected/total: %v/%v\n", api.Bot.Stats.ConnectedAccounts, api.Bot.Stats.TotalAccounts)
	message += fmt.Sprintf("Memory info:%s\n", GetMemStats())
	api.Bot.PostMessage(channel, message, api)
}
