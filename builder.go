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

func LoginWithToken(tokenType TokenType, token string) (*Session, error) {
	logger.Trace("LoginWithToken() called")

	if token == "" {
		return nil, errors.New("Configuration.Token cannot be empty")
	}

	session := &Session{tokenType: tokenType, token: token, objects: make(map[Snowflake]interface{})}

	gateway := gatewayGetResponse{}
	err := session.doRequest("GET", EndPointBotGateway().URL, nil, &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url + "?v=5&encoding=json"
	session.shards = make([]*shard, gateway.Shards)

	return session, nil
}
