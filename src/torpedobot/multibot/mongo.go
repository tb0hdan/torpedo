package multibot

import (
	"flag"
	"torpedobot/common"
	"torpedobot/common/database"
)

func (tb *TorpedoBot) ConfigureMongoDBPlugin() {
	tb.Config.MongoDBConnection = flag.String("mongo", "", "MongoDB server hostname")

}

func (tb *TorpedoBot) ParseMongoDBPlugin() {
	if *tb.Config.MongoDBConnection == "" {
		// try supplied one first
		*tb.Config.MongoDBConnection = common.GetStripEnv("MONGO")
		// docker...
		if *tb.Config.MongoDBConnection == "" {
			*tb.Config.MongoDBConnection = common.GetStripEnv("MONGO_PORT_27017_TCP_ADDR")
		}

	}
	tb.Database = database.New(*tb.Config.MongoDBConnection, "")
}
