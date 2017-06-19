package multibot

import (
	"flag"
	"torpedobot/common"
)

func (tb *TorpedoBot) ConfigureGooglePlugin() {
	tb.Config.GoogleWebAppKey = flag.String("google_webapp_key", "", "Google Data API Web Application Key")

}

func (tb *TorpedoBot) ParseGooglePlugin() {
	if *tb.Config.GoogleWebAppKey == "" {
		*tb.Config.GoogleWebAppKey = common.GetStripEnv("GOOGLE_WEBAPP_KEY")
	}
}
