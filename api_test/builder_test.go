package disgo_test

import (
	"flag"
	"os"
	"os/signal"
	"testing"

	"github.com/ikkerens/disgo"
	"github.com/slf4go/logger"
)

var (
	appQuit chan bool
	token   string
)

func init() {
	flag.StringVar(&token, "token", "", "Token for the bot")
	flag.Parse()

	appQuit = make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() { <-signalChan; appQuit <- true }()
}

func TestBuilder(t *testing.T) {
	discord, err := disgo.LoginWithToken(disgo.TypeBot, token)
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
		return
	}

	discord.RegisterEventHandler(onReady)

	err = discord.Connect()
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
		return
	}
	defer discord.Close()

	<-appQuit
}

func onReady(_ *disgo.Session, event disgo.ReadyEvent) {
	logger.Infof("onReady was called, logged in as %s with ID %d!", event.User.Username(), event.User.ID())
}
