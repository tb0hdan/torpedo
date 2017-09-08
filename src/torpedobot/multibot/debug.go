package multibot

import (
	"flag"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
)

var Debug *bool

func (tb *TorpedoBot) ConfigureDebug(cfg *torpedo_registry.ConfigStruct) {
	Debug = flag.Bool("debug", false, "Enable debug mode")
}

func bool2YN(in bool) (yn string) {
	yn = "no"
	if in {
		yn = "yes"
	}
	return
}

func (tb *TorpedoBot) ParseDebug(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("debug", bool2YN(*Debug))
	if cfg.GetConfig()["debug"] == "no" {
		cfg.SetConfig("debug", bool2YN(common.GetStripEnv("TORPEDO_DEBUG") != ""))
	}
}
