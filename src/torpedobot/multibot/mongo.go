package multibot

import (
	"flag"

	common "github.com/tb0hdan/torpedo_common"
	database "github.com/tb0hdan/torpedo_common/database"
	"github.com/tb0hdan/torpedo_registry"
)

var MongoDBConnection *string

func (tb *TorpedoBot) ConfigureMongoDBPlugin(cfg *torpedo_registry.ConfigStruct) {
	MongoDBConnection = flag.String("mongo", "", "MongoDB server hostname")
}

func (tb *TorpedoBot) ParseMongoDBPlugin(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("mongo", *MongoDBConnection)
	if cfg.GetConfig()["mongo"] == "" {
		// try supplied one first
		cfg.SetConfig("mongo", common.GetStripEnv("MONGO"))
		// docker...
		if cfg.GetConfig()["mongo"] == "" {
			cfg.SetConfig("mongo", common.GetStripEnv("MONGO_PORT_27017_TCP_ADDR"))
		}
		// fallback to localhost
		if cfg.GetConfig()["mongo"] == "" {
			cfg.SetConfig("mongo", "localhost")
		}

	}
	tb.Database = database.New(cfg.GetConfig()["mongo"], "")
}
