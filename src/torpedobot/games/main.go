package games

import (
	"log"
	"torpedobot/common"
	"os"
	"fmt"
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
	client = &Client{}
	client.logger = log.New(os.Stdout, fmt.Sprintf("%s: ", prefix), log.Lshortfile|log.LstdFlags)
	client.utils = &common.Utils{}
	client.utils.SetLoggerPrefix(prefix)
	client.StoreURL = store_url
	return
}
