package disgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
	"io/ioutil"
)

type Session struct {
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

func (s *Session) authorizationHeader() string {
	return s.tokenType.prefix + " " + s.token
}

func (s *Session) doHttpGet(url string, target interface{}) (err error) {
	response, err := s.doRequest("GET", url, nil)

	body := response.Body
	defer body.Close()
	if err = json.NewDecoder(body).Decode(target); err != nil {
		return err
	}

	return nil
}

func (s *Session) doHttpDelete(url string) error {
	_, err := s.doRequest("DELETE", url, nil)
	return err
}

func (s *Session) doHttpPost(url string, body interface{}) (err error) {
	jsonBody, err := json.Marshal(body)

	if err == nil {
		byteBuf := bytes.NewReader(jsonBody)
		_, err = s.doRequest("POST", url, byteBuf)
	}

	return err
}

func (s *Session) doRequest(method, url string, body io.Reader) (response *http.Response, err error) {
	path := strings.Replace(url, BaseUrl, "", 1)
	logger.Debugf("HTTP %s %s", method, path)

	var (
		req    *http.Request
		client = http.Client{
			Timeout: 10 * time.Second,
		}
	)

	if req, err = http.NewRequest(method, url, body); err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.authorizationHeader())
	req.Header.Add("User-Agent", "DiscordBot (https://github.com/ikkerens/disgo, 1.0.0)")

	if response, err = client.Do(req); err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 399 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Discord replied with status code %d: %s", response.StatusCode, string(body))
	}

	return response, nil
}
