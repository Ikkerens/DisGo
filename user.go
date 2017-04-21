package disgo

import "fmt"

func (s *User) AvatarURL() string {
	return fmt.Sprintf(EndPointUserAvatar, s.internal.ID, s.internal.AvatarHash)
}
