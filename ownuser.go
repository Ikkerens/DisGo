package disgo

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

type modifyCurrentUser struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func (s *Session) SetUsername(username string) (*User, error) {
	return s.modifyCurrentUser(modifyCurrentUser{Username: username})
}

func (s *Session) SetAvatar(imageMimeType string, reader io.Reader) (*User, error) {
	bytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	return s.modifyCurrentUser(modifyCurrentUser{Avatar: fmt.Sprintf("data:%s;base64,%s", imageMimeType, encoded)})
}

func (s *Session) modifyCurrentUser(modification modifyCurrentUser) (*User, error) {
	user := &User{}
	err := s.doHttpPatch(EndPointOwnUser(), &modification, user)
	if err != nil {
		return nil, err
	}
	user = objects.registerUser(user)
	if user.session == nil {
		user.session = s
	}
	return user, nil
}

type Status string

const (
	StatusOnline    Status = "online"
	StatusDND       Status = "dnd"
	StatusIdle      Status = "idle"
	StatusInvisible Status = "invisible"
)

func (s *Session) Status() Status {
	return s.status
}

func (s *Session) SetStatus(status Status) {
	s.SetStatusGame(status, s.game)
}

func (s *Session) SetGame(game *Game) {
	if s.status == "" {
		s.status = StatusOnline
	}

	s.SetStatusGame(s.status, game)
}

func (s *Session) SetStatusGame(status Status, game *Game) {
	var since uint64 = 0

	if status == StatusIdle {
		since = uint64(time.Now().Unix() * 1000)
	}

	s.status = status
	s.game = game

	for _, shard := range s.shards {
		shard.sendFrame(&gatewayFrame{opStatusUpdate, &statusPayload{
			Game:   game,
			Since:  since,
			Status: status,
			AFK:    status == StatusIdle,
		}}, false)
	}
}

func (s *shard) setSelfbotStatus() {
	s.sendFrame(&gatewayFrame{opStatusUpdate, &statusPayload{
		AFK:    true,
		Status: StatusInvisible,
	}}, false)
}

func (s *Session) GetDMChannel(userID Snowflake) (*Channel, error) {
	recipient := struct {
		RecipientID Snowflake `json:"recipient_id"`
	}{userID}

	channel := &Channel{}
	err := s.doHttpPost(EndPointDMChannels(), recipient, channel)
	if err != nil {
		return nil, err
	}
	channel = objects.registerChannel(channel)
	if channel.session == nil {
		channel.session = s
	}
	return channel, nil
}
