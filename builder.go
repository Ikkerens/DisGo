package disgo

import (
	"errors"
	"github.com/slf4go/logger"
)

type TokenType struct {
	prefix string
}

var (
	TypeBearer = TokenType{"Bearer"}
	TypeBot    = TokenType{"Bot"}
)

type gateWayResponse struct {
	Url    string `json:"url,omitempty"`
	Shards int    `json:"shards,omitempty"`
}

func WithToken(tokenType TokenType, token string) (*Session, error) {
	logger.Trace("WithToken() called")

	if token == "" {
		return nil, errors.New("Configuration.Token cannot be empty")
	}

	session := &Session{authorization: tokenType.prefix + " " + token}

	gateway := gateWayResponse{}
	err := session.doHttpGet(EndPointBotGateway, &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url
	session.shards = gateway.Shards

	return session, nil
}
