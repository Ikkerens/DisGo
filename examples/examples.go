package examples

import (
	"flag"
	"os"
	"os/signal"

	"github.com/slf4go/logger"
)

var (
	appQuit chan bool
	Token   string
)

func init() {
	logger.SetLevel(logger.LogAll)

	flag.StringVar(&Token, "token", "", "Token for the bot")
	flag.Parse()

	appQuit = make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		appQuit <- true

		go func() {
			<-signalChan
			os.Exit(1)
		}()
	}()
}

func WaitForExit() {
	<-appQuit
}
