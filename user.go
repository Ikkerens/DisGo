package disgo

import (
	"fmt"
	"time"
)

func (s *User) DiscordJoinDate() time.Time {
	return s.internal.ID.Timestamp()
}

func (s *User) AvatarURL() string {
	if s.internal.AvatarHash == "" {
		return ""
	}

	return fmt.Sprintf(EndPointUserAvatar, s.internal.ID, s.internal.AvatarHash)
}

func (s *User) Mention() string {
	return fmt.Sprintf("<@%s>", s.ID())
}
