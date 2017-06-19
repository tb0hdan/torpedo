package multibot

import (
	"flag"
	"torpedobot/common"
)

func (tb *TorpedoBot) ConfigurePinterestPlugin() {
	tb.Config.PinterestToken = *flag.String("pinterest_token", "", "Pinterest Client Token")

}

func (tb *TorpedoBot) RunPinterestPlugin() {
	if tb.Config.PinterestToken == "" {
		tb.Config.PinterestToken = common.GetStripEnv("PINTEREST")
	}
}
