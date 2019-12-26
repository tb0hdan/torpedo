package multibot

import (
	"fmt"
	"strconv"
	"time"

	"regexp"

	"gopkg.in/mgo.v2/bson"
)

type BlackListRuleItem struct {
	// sender, message
	Type string
	// mongo regexp
	Pattern string
	// match count
	Matches int64
}

func (tb *TorpedoBot) UpdateChannelTS(channel interface{}, message string) {
	var values []string
	chatID := fmt.Sprintf("%+v", channel)
	tb.logger.Printf("New channel")
	values = make([]string, 2)
	values[0] = fmt.Sprintf("%f", float64(time.Now().Unix()))
	values[1] = message
	tb.throttle.Set(chatID, values)
}

func (tb *TorpedoBot) GetChannelTSValues(channel interface{}) (values []string, ok bool) {
	chatID := fmt.Sprintf("%+v", channel)
	value, ok := tb.throttle.Get(chatID)
	values = value.([]string)
	return
}

func (tb *TorpedoBot) CheckMessageBlacklistOk(api *TorpedoBotAPI, message string) (status bool) {
	status = true
	session, collection, err := tb.Database.GetCollection("blackListItems")
	defer session.Close()
	if err != nil {
		tb.logger.Fatal("Could not connect to database: %+v\n", err)
	}
	result := []BlackListRuleItem{}
	err = collection.Find(bson.M{}).All(&result)
	for idx, filterItem := range result {
		tb.logger.Printf("Running filter check #%v\n", idx)
		switch filterItem.Type {
		case "message":
			r := regexp.MustCompile(filterItem.Pattern)
			if r.FindStringSubmatch(message) != nil {
				err = collection.Update(bson.M{"pattern": filterItem.Pattern}, bson.M{"matches": filterItem.Matches + 1,
					"type": filterItem.Type, "pattern": filterItem.Pattern})
				if err != nil {
					tb.logger.Fatalf("Could not update record %+v - %+v\n", filterItem, err)
				}
				status = false
				break
			}
		case "sender":
			r := regexp.MustCompile(filterItem.Pattern)
			if r.FindStringSubmatch(api.UserProfile.ID) != nil {
				status = false
				break
			}
		default:
			tb.logger.Fatalf("Got unknown filter type: `%s`\n", filterItem.Type)
		}
	}
	return
}

func (tb *TorpedoBot) NoSpam(api *TorpedoBotAPI, channel interface{}, message string) (status bool) {
	// empty messages are ok
	if message == "" {
		return true
	}
	// blacklisted sender/message check
	status = tb.CheckMessageBlacklistOk(api, message)
	if !status {
		tb.logger.Printf("Message blacklisted: %s", message)
		return
	}
	// message rate check
	values, ok := tb.GetChannelTSValues(channel)
	if !ok {
		tb.UpdateChannelTS(channel, message)
		status = true
	} else if val, err := strconv.ParseFloat(values[0], 64); err == nil && float64(time.Now().Unix()) > val+1 {
		tb.UpdateChannelTS(channel, message)
		status = true
	} else {
		tb.logger.Printf("Message rate exceeded for chat: %s", fmt.Sprintf("%+v", channel))
	}
	return
}
