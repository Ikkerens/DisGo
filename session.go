package disgo

import (
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

type Session struct {
	token     string
	tokenType string
	wsUrl     string

	shards       []*shard
	shuttingDown bool

	objects map[Snowflake]interface{}
}

func BuildWithBotToken(token string) (*Session, error) {
	logger.Trace("BuildWithBotToken() called")

	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	session := &Session{tokenType: "Bot", token: token, objects: make(map[Snowflake]interface{})}

	gateway := gatewayGetResponse{}
	_, err := session.doRequest("GET", EndPointBotGateway().URL, nil, &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url + "?v=5&encoding=json"
	session.shards = make([]*shard, gateway.Shards)

	return session, nil
}

func (s *Session) Connect() error {
	for i := 0; i < cap(s.shards); i++ {
		shard, err := connectShard(s, i)
		if err != nil {
			s.closeShards(websocket.CloseAbnormalClosure, fmt.Sprintf("Error occured on shard [%d/%d]", i, cap(s.shards)))
			return err
		}

		s.shards[i] = shard

		if (i + 1) != cap(s.shards) {
			time.Sleep(5 * time.Second)
		}
	}

	return nil
}

func (s *Session) Close() {
	s.shuttingDown = true
	s.closeShards(websocket.CloseNormalClosure, "")
}

func (s *Session) closeShards(code int, text string) {
	for _, sh := range s.shards {
		if sh != nil {
			sh.disconnect(code, text)
		}
	}
}
