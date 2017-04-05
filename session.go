package disgo

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/slf4go/logger"
)

type Session struct {
	authorization string
	wsUrl         string
	shards        int
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

	req.Header.Add("Authorization", s.authorization)
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
