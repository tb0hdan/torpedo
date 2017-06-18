package games

import (
	"log"
	"torpedobot/common"
)

type Client struct {
	StoreURL string
	logger *log.Logger
	utils *common.Utils
}

type GameItem struct {
	GameURL string
	GameThumbnail string
	CurrentPrice float64
	RegularPrice float64
	DiscountPercentage int
	Platforms []string
}

func (gi *GameItem) IsComplete() (result bool) {
	if gi.GameURL != "" && gi.GameThumbnail != "" {
		result = true
	}
	return
}

func  NewClient(prefix, store_url string) (client *Client){
	cu := &common.Utils{}
	client = &Client{}
	client.logger = cu.NewLog(prefix)
	client.utils = &common.Utils{}
	client.utils.SetLoggerPrefix(prefix)
	client.StoreURL = store_url
	return
}
