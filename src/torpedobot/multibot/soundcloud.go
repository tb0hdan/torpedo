package multibot

import (
	"flag"
	"torpedobot/common"
)

func (tb *TorpedoBot) ConfigureSoundCloudPlugin() {
	tb.Config.SoundCloudClientID = flag.String("soundcloud_id", "", "SoundCloud client ID")

}

func (tb *TorpedoBot) ParseSoundCloudPlugin() {
	if *tb.Config.SoundCloudClientID == "" {
		*tb.Config.SoundCloudClientID = common.GetStripEnv("SOUNDCLOUD_ID")
	}
}
