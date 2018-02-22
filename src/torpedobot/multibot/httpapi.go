package multibot

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ant0ine/go-json-rest/rest"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
)

var APIADDR *string

func (tb *TorpedoBot) ConfigureHTTPAPI(cfg *torpedo_registry.ConfigStruct) {
	APIADDR = flag.String("apiaddr", "", "Listen on this address for incoming HTTP API Server connections. Example: :8080")
}

func (tb *TorpedoBot) ParseHTTPAPI(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("apiaddr", *APIADDR)
	if cfg.GetConfig()["apiaddr"] == "" {
		// try supplied one first
		cfg.SetConfig("apiaddr", common.GetStripEnv("APIADDR"))
	}
}

func (tb *TorpedoBot) RunHTTPAPI() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
		w.WriteJson(map[string]string{"Body": "Hello World!"})
	}))

	apiaddr := torpedo_registry.Config.GetConfig()["apiaddr"]

	if apiaddr != "" {
		log.Fatal(http.ListenAndServe(apiaddr, api.MakeHandler()))
	}
}
