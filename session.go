package disgo

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

const gatewayVersion = "6"

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

func NewSelfBot(token string) (*Session, error) {
	logger.Trace("NewSelfBot() called")

	if token == "" {
		panic("token cannot be empty")
	}

	session := &Session{tokenType: "", token: token, rateLimitBuckets: make(map[string]*rateBucket)}
	registerInternalEvents(session)

	gateway := gatewayGetResponse{}
	err := session.doHttpGet(EndPointGateway(), &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url + "?v=" + gatewayVersion + "&encoding=json"
	session.shards = make([]*shard, 1)

	return session, nil
}

func NewBot(token string) (*Session, error) {
	logger.Trace("NewBot() called")

	if token == "" {
		panic("token cannot be empty")
	}

	session := &Session{tokenType: "Bot ", token: token, rateLimitBuckets: make(map[string]*rateBucket)}
	registerInternalEvents(session)

	gateway := gatewayGetResponse{}
	err := session.doHttpGet(EndPointBotGateway(), &gateway)
	if err != nil {
		return nil, err
	}

	session.wsUrl = gateway.Url + "?v=" + gatewayVersion + "&encoding=json"
	session.shards = make([]*shard, gateway.Shards)

	return session, nil
}

func registerInternalEvents(session *Session) {
	session.registerEventHandler(onChannelCreate, false)
	session.registerEventHandler(onChannelDelete, false)
	session.registerEventHandler(onGuildMemberUpdate, false)
	session.registerEventHandler(onGuildMemberAdd, false)
}

func (s *Session) Connect() error {
	for i := 0; i < cap(s.shards); i++ {
		shard, err := newShard(s, i)
		if err != nil {
			s.closeShards(websocket.CloseGoingAway, fmt.Sprintf("Error occurred on shard [%d/%d]", i, cap(s.shards)))
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
