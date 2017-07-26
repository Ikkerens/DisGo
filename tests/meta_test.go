package tests

import (
	"flag"
	"fmt"
	"github.com/ikkerens/disgo"
	"github.com/slf4go/logger"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	discord *disgo.Session
	guildID disgo.Snowflake

	wsChannel = make(chan idObject)
)

func testMain(m *testing.M, eventSetter func()) {
	fmt.Println("=== SETUP   DisGo public api tests")
	logger.SetLevel(logger.LogAll)
	rand.Seed(time.Now().Unix())

	// Get the token from args
	var token string
	flag.StringVar(&token, "token", "", "Token for the automated tests.")
	flag.Parse()

	// Create a new REST instance
	var err error
	discord, err = disgo.NewBot(token)
	if err != nil {
		logger.ErrorE(err)
		os.Exit(1)
		return
	}

	// Prepare fetching the testing guild id
	guildIDChan := make(chan disgo.Snowflake)
	discord.RegisterEventHandler(func(_ *disgo.Session, event disgo.GuildCreateEvent) {
		guildIDChan <- event.ID()
	})
	eventSetter()

	// Connect to the websocket
	err = discord.Connect()
	if err != nil {
		os.Exit(1)
		return
	}

	// Set up the guild ID for the tests
	guildID = <-guildIDChan
	close(guildIDChan)

	// Run the tests
	testResult := m.Run()

	fmt.Println("=== TEARDOWN   Disgo public api tests")
	// Close the connection gracefully
	discord.Close()
	// Return the test re sults
	os.Exit(testResult)
}

func errorCheck(t *testing.T, err error) {
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
	}
}

type idObject interface {
	ID() disgo.Snowflake
}

func waitForMirror(t *testing.T, object idObject) {
	timer := time.After(5 * time.Second)
loop:
	for {
		select {
		case wsObject := <-wsChannel:
			if wsObject.ID() == object.ID() {
				break loop
			}
		case <-timer:
			logger.Error("No copy of object received on the websocket.")
			t.FailNow()
		}
	}
}
