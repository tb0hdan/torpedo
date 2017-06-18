package bashorg

import (
	"bytes"
	"log"

	"torpedobot/common"

	"golang.org/x/net/html"
)


type BashOrg struct {
	logger *log.Logger
	utils  *common.Utils
}


func (bo *BashOrg) Get_bashorg_html(url string) (result *html.Node) {
	res, err := bo.utils.GetURLBytes(url)
	if err != nil {
		return
	}
	result, err = html.Parse(bytes.NewReader(res))
	if err != nil {
		bo.logger.Printf("Error %s parsing html", err)
	}
	return
}


func NewClient() (client *BashOrg) {
	client = &BashOrg{}
	client.utils = &common.Utils{}
	client.logger = client.utils.SetLoggerPrefix("bashorg-plugin")
	return
}
