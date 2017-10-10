package disgo

import "strconv"

func (s *Guild) BuildChannel(name string) *ChannelBuilder {
	return s.session.BuildChannel(s.internal.ID, name)
}

func (s *Guild) Role(id Snowflake) (*Role, bool) {
	for _, role := range s.internal.Roles {
		if role.internal.ID == id {
			return role, true
		}
	}

	return nil, false
}

func (s *Session) AddGuildMemberRole(guildID, userID, roleID Snowflake) error {
	return s.doHttpPut(EndPointGuildMemberRoles(guildID, userID, roleID), nil)
}

func (s *Session) RemoveGuildMemberRole(guildID, userID, roleID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMemberRoles(guildID, userID, roleID), nil)
}

func (s *Guild) GetRoleUsers(roleID Snowflake) []*User {
	users := make([]*User, 0)

	for _, user := range s.internal.Members {
		if SnowflakeInSlice(roleID, user.internal.RolesIDs) {
			users = append(users, user.internal.User)
		}
	}

	return users
}

func (s *Guild) GetUserRoles(userID Snowflake) ([]Snowflake, bool) {
	membership, exists := s.GetUserMembership(userID)
	if exists {
		return membership.RolesIDs(), true
	}

	return nil, false
}

func (s *Guild) GetUserMembership(userID Snowflake) (*GuildMember, bool) {
	for _, membership := range s.internal.Members {
		if membership.User().ID() == userID {
			return membership, true
		}
	}

	return nil, false
}

func (s *Guild) GetUserColor(userID Snowflake) (int, bool) {
	roles, inGuild := s.GetUserRoles(userID)
	if !inGuild {
		return 0, false
	}

	var (
		highest *Role
	)
	for _, role := range s.internal.Roles {
		if SnowflakeInSlice(role.internal.ID, roles) && (highest == nil || role.internal.Position > highest.internal.Position) && role.internal.Color != 0 {
			highest = role
		}
	}

	if highest != nil {
		return highest.internal.Color, true
	}

	return 0, false
}

type updateGuildMember struct {
	Nick      string       `json:"nick,omitempty"`
	Roles     *[]Snowflake `json:"roles,omitempty"`
	Mute      *bool        `json:"mute,omitempty"`
	Deaf      *bool        `json:"deaf,omitempty"`
	ChannelID *Snowflake   `json:"channel_id,omitempty"`
}

func (s *Session) SetUserNick(guildID, userID Snowflake, nick string) error {
	return s.doHttpPatch(EndPointGuildMember(guildID, userID), updateGuildMember{Nick: nick}, nil)
}

func (s *Session) KickUser(guildID, userID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMember(guildID, userID), nil)
}

func (s *Guild) KickUser(userID Snowflake) error {
	return s.session.KickUser(s.internal.ID, userID)
}

func (s *Session) BanUser(guildID, userID Snowflake, deleteMessageDays int) error {
	endPoint := EndPointGuildMemberBan(guildID, userID)
	endPoint.Url += "?delete-message-days=" + strconv.FormatInt(int64(deleteMessageDays), 10)

	return s.doHttpPut(endPoint, nil)
}

func (s *Guild) BanUser(userID Snowflake, deleteMessageDays int) error {
	return s.session.BanUser(s.internal.ID, userID, deleteMessageDays)
}

func (s *Session) UnbanUser(guildID, userID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMemberBan(guildID, userID), nil)
}

func (s *Guild) UnbanUser(userID Snowflake) error {
	return s.session.UnbanUser(s.internal.ID, userID)
}
