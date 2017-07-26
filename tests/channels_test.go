package tests

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/ikkerens/disgo"
)

var (
	testChannel *disgo.Channel
	testMessage *disgo.Message
)

// Set up event handlers for tests
func TestMain(m *testing.M) {
	testMain(m, func() {
		discord.RegisterEventHandler(func(_ *disgo.Session, event disgo.ChannelCreateEvent) {
			wsChannel <- event.Channel
		})
		discord.RegisterEventHandler(func(_ *disgo.Session, event disgo.MessageCreateEvent) {
			wsChannel <- event.Message
		})
	})
}

func TestCreateChannel(t *testing.T) {
	channelName := "test-" + strconv.FormatInt(rand.Int63n(1000000), 16)
	channel, err := discord.BuildChannel(guildID, channelName).Create()
	errorCheck(t, err)

	waitForMirror(t, channel)
	testChannel = channel
}

func TestSendMessage(t *testing.T) {
	// This test depends on the result of a previous test.
	if testChannel == nil {
		t.SkipNow()
	}

	message, err := discord.SendMessage(testChannel.ID(), "Hi there! This is a test!")
	errorCheck(t, err)

	waitForMirror(t, message)
	testMessage = message
}

func TestDeleteMessage(t *testing.T) {
	// This test depends on the result of a previous test.
	if testMessage == nil {
		t.SkipNow()
	}

	err := testMessage.Delete()
	errorCheck(t, err)
}

func TestDeleteChannel(t *testing.T) {
	// This test depends on the result of a previous test.
	if testChannel == nil {
		t.SkipNow()
	}

	err := testChannel.Delete()
	errorCheck(t, err)
}
