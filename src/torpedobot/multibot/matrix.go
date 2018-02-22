package multibot

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"flag"

	common "github.com/tb0hdan/torpedo_common"

	"github.com/matrix-org/gomatrix"
	"github.com/tb0hdan/torpedo_registry"
)

var MatrixAPIKey *string

func HandleMatrixMessage(channel interface{}, message string, tba *TorpedoBotAPI, richmsgs []torpedo_registry.RichMessage) {
	switch api := tba.API.(type) {
	case *gomatrix.Client:
		if len(richmsgs) > 0 && !richmsgs[0].IsEmpty() {
			msg, url := richmsgs[0].ToGenericAttachment()
			api.SendImage(channel.(string), msg, url)
		} else {
			api.SendText(channel.(string), message)
		}
	}
}

func (tb *TorpedoBot) ConfigureMatrixBot(cfg *torpedo_registry.ConfigStruct) {
	MatrixAPIKey = flag.String("matrix", "", "Matrix.org creds: ID:AccessToken,")
}

func (tb *TorpedoBot) ParseMatrixBot(cfg *torpedo_registry.ConfigStruct) {
	cfg.SetConfig("matrixapikey", *MatrixAPIKey)
	if cfg.GetConfig()["matrixapikey"] == "" {
		cfg.SetConfig("matrixapikey", common.GetStripEnv("MATRIX"))
	}
}

func (tb *TorpedoBot) RunMatrixBot(apiKey, cmd_prefix string) {
	account := &torpedo_registry.Account{
		APIKey:        apiKey,
		CommandPrefix: cmd_prefix,
	}
	torpedo_registry.Accounts.AppendAccounts(account)
	tb.RunMatrixBotAccount(account)
}

func (tb *TorpedoBot) RunMatrixBotAccount(account *torpedo_registry.Account) {
	tb.Stats.ConnectedAccounts += 1
	account.Connection.ReconnectCount += 1

	cu := &common.Utils{}

	logger := cu.NewLog("matrix-bot")

	clientID := fmt.Sprintf("@%s:matrix.org", strings.Split(account.APIKey, ":")[0])
	cli, _ := gomatrix.NewClient("https://matrix.org", clientID, strings.Split(account.APIKey, ":")[1])
	// anything which implements the Storer interface
	customStore := gomatrix.NewInMemoryStore()
	cli.Store = customStore

	// anything which implements the Syncer interface
	customSyncer := gomatrix.NewDefaultSyncer(clientID, customStore)
	cli.Syncer = customSyncer

	// any http.Client
	cli.Client = http.DefaultClient

	account.API = cli
	syncer := cli.Syncer.(*gomatrix.DefaultSyncer)
	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		logger.Printf("Message: %+v\n", ev)
		if ev.Sender != clientID {

			botApi := &TorpedoBotAPI{}
			botApi.API = cli
			botApi.Bot = tb
			botApi.CommandPrefix = account.CommandPrefix
			botApi.UserProfile = &torpedo_registry.UserProfile{ID: ev.Sender}
			botApi.Me = clientID

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

	tb.RegisteredProtocols["*gomatrix.Client"] = HandleMatrixMessage

	for {
		if err := cli.Sync(); err != nil {
			logger.Printf("Sync() failed with: %+v\n", err)
		}
		// Optional: Wait a period of time before trying to sync again.
		time.Sleep(10 * time.Second)
	}
	tb.Stats.ConnectedAccounts -= 1
}
