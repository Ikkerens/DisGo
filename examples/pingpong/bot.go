package main

import (
	"github.com/ikkerens/disgo"
	"github.com/ikkerens/disgo/examples"
	"github.com/slf4go/logger"
)

func main() {
	discord, err := disgo.NewBot(examples.Token)
	if err != nil {
		logger.ErrorE(err)
		return
	}

	discord.RegisterEventHandler(onMessage)

	err = discord.Connect()
	if err != nil {
		logger.ErrorE(err)
		return
	}
	defer discord.Close()

	examples.WaitForExit()
}

func onMessage(_ *disgo.Session, event disgo.MessageCreateEvent) {
	if event.Content() == "!ping" {
		event.Reply("Pong!")
	}
}
