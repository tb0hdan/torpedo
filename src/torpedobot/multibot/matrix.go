package multibot

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"log"
	"os"

	"github.com/matrix-org/gomatrix"
)

func (tb *TorpedoBot) RunMatrixBot(apiKey, cmd_prefix string) {
	logger := log.New(os.Stdout, "matrix-bot: ", log.Lshortfile|log.LstdFlags)

	clientID := fmt.Sprintf("@%s:matrix.org", strings.Split(apiKey, ":")[0])
	cli, _ := gomatrix.NewClient("https://matrix.org", clientID, strings.Split(apiKey, ":")[1])
	// anything which implements the Storer interface
	customStore := gomatrix.NewInMemoryStore()
	cli.Store = customStore

	// anything which implements the Syncer interface
	customSyncer := gomatrix.NewDefaultSyncer(clientID, customStore)
	cli.Syncer = customSyncer

	// any http.Client
	cli.Client = http.DefaultClient

	syncer := cli.Syncer.(*gomatrix.DefaultSyncer)
	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		logger.Printf("Message: %+v\n", ev)
		if ev.Sender != clientID {

			botApi := &TorpedoBotAPI{}
			botApi.API = cli
			botApi.Bot = tb
			botApi.CommandPrefix = cmd_prefix

			msg, _ := ev.Body()
			go tb.processChannelEvent(botApi, ev.RoomID, msg)
		}

	})
	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		logger.Printf("Room: %+v\n", ev)
		servername := ""
		if !strings.HasSuffix(ev.RoomID, ":matrix.org") {
			servername = strings.Split(ev.RoomID, ":")[1]
		}
		cli.JoinRoom(ev.RoomID, servername, "")
	})

	logger.Printf("Starting Matrix.Org bot...")
	for {
		if err := cli.Sync(); err != nil {
			logger.Printf("Sync() failed with: %+v\n", err)
		}
		// Optional: Wait a period of time before trying to sync again.
		time.Sleep(10 * time.Second)
	}
}
