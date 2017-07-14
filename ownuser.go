package disgo

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

type modifyCurrentUser struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func (s *Session) SetUsername(username string) (*User, error) {
	return s.modifyCurrentUser(modifyCurrentUser{Username: username})
}

func (s *Session) SetAvatar(imageMimeType string, reader io.Reader) (*User, error) {
	bytes, _ := ioutil.ReadAll(reader)
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
