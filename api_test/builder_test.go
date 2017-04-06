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
	session, err := disgo.BuildWithToken(disgo.TypeBot, token)
	if err != nil {
		logger.ErrorE(err)
		t.FailNow()
		return
	}
	session.Open()
	session.Close()
}

func TestInvalidSessionCheck(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.FailNow()
		}
	}()

	session := disgo.Session{}
	session.Close()
}
