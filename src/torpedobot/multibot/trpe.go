package multibot

import (
	"net/url"

	"flag"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
)

var TRPEURL *string

type TRPEResponse struct {
}

func (tb *TorpedoBot) ConfigureTRPE(cfg *torpedo_registry.ConfigStruct) {
	TRPEURL = flag.String("trpe_host", "", "TRPE URL (disabled if unset)")

}

func (tb *TorpedoBot) ParseTRPE(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("trpe_host", *TRPEURL)
	if cfg.GetConfig()["trpe_host"] == "" {
		cfg.SetConfig("trpe_host", common.GetStripEnv("TRPE_HOST"))
	}
}

func (tb *TorpedoBot) processViaTRPE(channel interface{}, incoming_message, command_prefix, host string) (err error, result string) {
	cu := common.Utils{}
	response := &TRPEResponse{}
	err = cu.PostURLFormUnmarshal(host, url.Values{"key": {"Value"}, "id": {"123"}}, response)
	return
}
