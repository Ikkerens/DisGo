package disgo_test

import (
	"flag"
	"testing"

	"github.com/ikkerens/disgo"
	"github.com/slf4go/logger"
)

var token string

func init() {
	flag.StringVar(&token, "token", "", "Token for the bot")
	flag.Parse()
}

func TestBuilder(t *testing.T) {
	session, err := disgo.WithToken(disgo.TypeBot, token)
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
		return
	}
	session.Open()
}
