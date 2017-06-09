package main

import (
	"fmt"
	"runtime"
	"time"
	"torpedobot/multibot"
)

func FormatTDiff(ts int64) (int64, int64, int64, int64) {
	m, s := ts/60, ts%60
	h, m := m/60, m%60
	d, h := h/24, h%24
	return d, h, m, s
}

func CalculateMessageRate(tdiff, msgcount int64) (result string) {
	var value int64
	pairs := make(map[string]int64)
	pairs["s"] = 0
	pairs["m"] = 60
	pairs["h"] = 3600
	pairs["d"] = 86400
	if msgcount == 0 {
		result = "0/s"
		return
	}
	// BUGGY!
	for key := range pairs {
		if pairs[key] == 0 {
			value = msgcount / tdiff
		} else {
			value = msgcount / (tdiff / pairs[key])
		}
		if value > 0 {
			result = fmt.Sprintf("%v/%s", value, key)
			break
		}
	}
	return
}

func GetMemStats() (result string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	result = fmt.Sprintf("\n\tAlloc = %v KB\n\tSys = %v KB\n\tNumGC = %v\n\n", m.Alloc/1024, m.Sys/1024, m.NumGC)
	return
}

func StatsProcessMessage(api *multibot.TorpedoBotAPI, channel interface{}, incoming_message string) {
	ts := int64(time.Now().Unix()) - api.Bot.Stats.StartTimestamp
	d, h, m, s := FormatTDiff(ts)
	message := fmt.Sprintf("Uptime: %v day(s) %v hour(s) %v minute(s) %v second(s)\n", d, h, m, s)
	message += fmt.Sprintf("Processed messages: %v\n", api.Bot.Stats.ProcessedMessages)
	//message += fmt.Sprintf("Message rate: %s\n", CalculateMessageRate(ts, api.Bot.Stats.ProcessedMessages))
	message += fmt.Sprintf("Accounts connected/total: %v/%v\n", api.Bot.Stats.ConnectedAccounts, api.Bot.Stats.TotalAccounts)
	message += fmt.Sprintf("Memory info:%s\n", GetMemStats())
	api.Bot.PostMessage(channel, message, api)
}
