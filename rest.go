package disgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/slf4go/logger"
)

func (s *Session) doHttpGet(endPoint EndPoint, target interface{}) (err error) {
	err = s.doRequest("GET", endPoint(false), nil, target)
	return err
}

func (s *Session) doHttpDelete(endPoint EndPoint, target interface{}) error {
	err := s.doRequest("DELETE", endPoint(false), nil, target)
	return err
}

func (s *Session) doHttpPost(endPoint EndPoint, body, target interface{}) (err error) {
	jsonBody, err := json.Marshal(body)

	if err == nil {
		byteBuf := bytes.NewReader(jsonBody)
		err = s.doRequest("POST", endPoint(false), byteBuf, target)
	}

	return err
}

func (s *Session) doRequest(method, url string, body io.Reader, target interface{}) (err error) {
	logger.Debugf("HTTP %s %s", method, strings.Replace(url, BaseUrl, "", 1))

	var (
		req      *http.Request
		response *http.Response
		client   = http.Client{
			Timeout: 10 * time.Second,
		}
	)

	if req, err = http.NewRequest(method, url, body); err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.tokenType.prefix+" "+s.token)
	req.Header.Add("User-Agent", "DiscordBot (https://github.com/ikkerens/disgo, 1.0.0)")

	if response, err = client.Do(req); err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 399 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Discord replied with status code %d: %s", response.StatusCode, string(body))
	} else if target != nil {
		body := response.Body
		defer body.Close()
		if err = json.NewDecoder(body).Decode(target); err != nil {
			return err
		}
	}

	return nil
}
