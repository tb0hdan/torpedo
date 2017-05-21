package multibot

import (
	"time"
	"fmt"
	"strconv"
)

func (tb *TorpedoBot) UpdateChannelTS(channel interface{}, message string) {
	var values []string
	chatID := fmt.Sprintf("%+v", channel)
	tb.logger.Printf("New channel")
	values = make([]string, 2)
	values[0] = fmt.Sprintf("%f",  float64(time.Now().Unix()))
	values[1] = message
	tb.throttle.Set(chatID, values)
}

func (tb *TorpedoBot) GetChannelTSValues(channel interface{}) (values []string, ok bool){
	chatID := fmt.Sprintf("%+v", channel)
	values, ok = tb.throttle.Get(chatID)
	return
}

func (tb *TorpedoBot) NoSpam(channel interface{}, message string) (status bool) {
	// Just message rate check so far
	values, ok := tb.GetChannelTSValues(channel)
	if ! ok {
		tb.UpdateChannelTS(channel, message)
		status = true
	} else if val, err := strconv.ParseFloat(values[0], 64); err == nil && float64(time.Now().Unix()) > val + 1 {
		tb.UpdateChannelTS(channel, message)
		status = true
	} else {
		tb.logger.Printf("Message rate exceeded for chat: %s", fmt.Sprintf("%+v", channel))
	}
	return
}

