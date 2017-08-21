package multibot

import (
	"fmt"
	"time"

	"github.com/tb0hdan/torpedo_registry"
)

func (tb *TorpedoBot) StoreMessageHistory(api *TorpedoBotAPI, channel interface{}, incoming_message string) {
	session, collection, err := tb.Database.GetCollection("chatHistory")
	defer session.Close()
	if err != nil {
		tb.logger.Fatal("Could not connect to database: %+v\n", err)
	}
	err = collection.Insert(&torpedo_registry.MessageHistoryItem{Timestamp: int64(time.Now().Unix()), Channel: fmt.Sprintf("%v", channel),
		Sender: api.UserProfile.ID, Message: incoming_message, Nick: api.UserProfile.Nick})
	if err != nil {
		tb.logger.Fatal(err)
	}
	tb.logger.Printf("%v - %s - %s\n", channel, api.UserProfile.ID, incoming_message)
	return
}
