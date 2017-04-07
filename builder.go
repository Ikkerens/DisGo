package disgo

import (
	"errors"
	"time"

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

	session := &Session{valid: true, tokenType: tokenType, token: token}

	gateway := gatewayGetResponse{}
	err := session.doHttpGet(EndPointBotGateway, &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url + "?v=5&encoding=json"
	session.shards = make([]*Shard, gateway.Shards)

	for i := 0; i < gateway.Shards; i++ {
		shard, err := connectShard(session, i)
		if err != nil {
			return nil, err
		}

		session.shards[i] = shard

		if (i + 1) != gateway.Shards {
			time.Sleep(5 * time.Second)
		}
	}

	return session, nil
}
