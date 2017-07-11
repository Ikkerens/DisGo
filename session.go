package disgo

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

type Session struct {
	token     string
	tokenType string
	wsUrl     string

	rateLimitBuckets map[string]*rateBucket
	globalRateLimit  sync.Mutex
	globalReset      time.Time

	shards       []*shard
	shuttingDown bool
	stateLock    sync.RWMutex
}

func BuildSelfbotWithToken(token string) (*Session, error) {
	logger.Trace("BuildSelfbotWithToken() called")

	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	session := &Session{tokenType: "", token: token, rateLimitBuckets: make(map[string]*rateBucket)}

	// Internal event handlers
	session.registerEventHandler(onGuildMemberUpdate, false)
	session.registerEventHandler(onGuildMemberAdd, false)

	gateway := gatewayGetResponse{}
	_, err := session.doRequest("GET", EndPointGateway().Url, nil, &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url + "?v=5&encoding=json"
	session.shards = make([]*shard, 1)

	return session, nil
}

func BuildWithBotToken(token string) (*Session, error) {
	logger.Trace("BuildWithBotToken() called")

	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	session := &Session{tokenType: "Bot ", token: token, rateLimitBuckets: make(map[string]*rateBucket)}

	// Internal event handlers
	session.registerEventHandler(onGuildMemberUpdate, false)
	session.registerEventHandler(onGuildMemberAdd, false)

	gateway := gatewayGetResponse{}
	_, err := session.doRequest("GET", EndPointBotGateway().Url, nil, &gateway)
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
			s.closeShards(websocket.CloseAbnormalClosure, fmt.Sprintf("Error occurred on shard [%d/%d]", i, cap(s.shards)))
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
	s.stateLock.Lock()
	s.shuttingDown = true
	s.stateLock.Unlock()
	s.closeShards(websocket.CloseNormalClosure, "")
}

func (s *Session) isShuttingDown() bool {
	s.stateLock.RLock()
	defer s.stateLock.RUnlock()

	return s.shuttingDown
}

func (s *Session) closeShards(code int, text string) {
	for _, sh := range s.shards {
		if sh != nil {
			sh.disconnect(code, text)
		}
	}
}
