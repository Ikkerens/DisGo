package disgo

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	token     string
	tokenType TokenType
	wsUrl     string

	shards       []*shard
	shuttingDown bool

	objects map[Snowflake]interface{}
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
