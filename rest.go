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
	err = s.rateLimit(endPoint, func() (*http.Response, error) {
		return s.doRequest("GET", endPoint.URL, nil, target)
	})
	return
}

func (s *Session) doHttpDelete(endPoint EndPoint, target interface{}) (err error) {
	err = s.rateLimit(endPoint, func() (*http.Response, error) {
		return s.doRequest("DELETE", endPoint.URL, nil, target)
	})
	return
}

func (s *Session) doHttpPost(endPoint EndPoint, body, target interface{}) (err error) {
	jsonBody, err := json.Marshal(body)

	if err == nil {
		byteBuf := bytes.NewReader(jsonBody)
		err = s.rateLimit(endPoint, func() (*http.Response, error) {
			return s.doRequest("POST", endPoint.URL, byteBuf, target)
		})
	}

	return
}

func (s *Session) doRequest(method, url string, body io.Reader, target interface{}) (response *http.Response, err error) {
	logger.Debugf("HTTP %s %s", method, strings.Replace(url, BaseUrl, "", 1))

	var (
		req    *http.Request
		client = http.Client{
			Timeout: 10 * time.Second,
		}
	)

	if req, err = http.NewRequest(method, url, body); err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.tokenType+" "+s.token)
	req.Header.Add("User-Agent", "DiscordBot (https://github.com/ikkerens/disgo, 1.0.0)")

	if response, err = client.Do(req); err != nil {
		return
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		fallthrough
	case 201:
		if target != nil {
			body := response.Body
			defer body.Close()
			if err = json.NewDecoder(body).Decode(target); err != nil {
				return
			}
		}
	case 204:
		fallthrough
	case 304:
		return
	default:
		var bodyBuf []byte
		bodyBuf, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
		return response, fmt.Errorf("Discord replied with status code %d: %s", response.StatusCode, string(bodyBuf))
	}

	return
}
