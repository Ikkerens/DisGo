package disgo

import "fmt"

func (s *User) AvatarURL() string {
	if s.internal.AvatarHash == "" {
		return ""
	}

	return fmt.Sprintf(EndPointUserAvatar, s.internal.ID, s.internal.AvatarHash)
}
