package disgo_test

import (
	"flag"
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
}

func TestMessageCreateDelete(t *testing.T) {
	discord, err := disgo.BuildWithBotToken(token)
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
		return
	}

	discord.RegisterEventHandler(onGuildCreate)
	discord.RegisterEventHandler(onMessage)

	err = discord.Connect()
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
		return
	}

	defer discord.Close()

	<-appQuit
}

func onGuildCreate(s *disgo.Session, event disgo.GuildCreateEvent) {
	for _, channel := range event.Channels() {
		if channel.Type() == "text" && channel.Name() == "bottest" {
			_, err := s.SendMessage(channel.ID(), "I am going to delete this message!")
			if err != nil {
				logger.ErrorE(err)
			}
		}
	}
}

func onMessage(s *disgo.Session, event disgo.MessageCreateEvent) {
	logger.Infof("User %s posted: %s", event.Author().Username(), event.Content())
	event.Message.Delete()
	appQuit <- true
}
