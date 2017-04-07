package disgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/slf4go/logger"
)

type Session struct {
	valid bool

	token     string
	tokenType TokenType
	wsUrl     string
	shards    []*Shard
}

func (s *Session) Close() {
	for _, sh := range s.shards {
		if  sh != nil {
			sh.disconnect()
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

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return err
	}

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	logger.Debugf("doHttpGet(%s) response: %+v", path, target)

	return nil
}
