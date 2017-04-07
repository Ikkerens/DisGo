package disgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

type Session struct {
	valid bool

	token     string
	tokenType TokenType
	wsUrl     string

	shards       []*Shard
	shuttingDown bool
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

func (s *Session) checkValid() {
	if !s.valid {
		caller, _, _, _ := runtime.Caller(1)
		panic(fmt.Sprintf("%s() requires a valid Session", runtime.FuncForPC(caller).Name()))
	}
}

func (s *Session) authorizationHeader() string {
	return s.tokenType.prefix + " " + s.token
}

func (s *Session) doHttpGet(url string, target interface{}) error {
	path := strings.Replace(url, BaseUrl, "", 1)
	logger.Tracef("doHttpGet(%s) called", path)

	var (
		req *http.Request
		err error
	)

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return err
	}

	req.Header.Add("Authorization", s.authorizationHeader())
	req.Header.Add("User-Agent", "DiscordBot (https://github.com/ikkerens/disgo, 1.0.0)")

	var (
		client = http.Client{
			Timeout: 10 * time.Second,
		}
		resp *http.Response
	)
	if resp, err = client.Do(req); err != nil {
		return err
	}

	body := resp.Body
	defer body.Close()
	if err = json.NewDecoder(body).Decode(target); err != nil {
		return err
	}

	logger.Debugf("doHttpGet(%s) response: %+v", path, target)

	return nil
}
