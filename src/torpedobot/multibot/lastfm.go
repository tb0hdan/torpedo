package multibot

import (
	"flag"
	"torpedobot/common"
)

func (tb *TorpedoBot) ConfigureLastFmPlugin() {
	tb.Config.LastFmKey = *flag.String("lastfm_key", "", "Last.FM API Key")
	tb.Config.LastFmSecret = *flag.String("lastfm_secret", "", "Last.FM API Secret")

}

func (tb *TorpedoBot) RunLastFmPlugin() {
	if tb.Config.LastFmKey == "" {
		tb.Config.LastFmKey = common.GetStripEnv("LASTFM_KEY")
	}
	if tb.Config.LastFmSecret == "" {
		tb.Config.LastFmSecret = common.GetStripEnv("LASTFM_SECRET")
	}
}
