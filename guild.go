package disgo

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

func (s *Guild) GetUserRoles(userID Snowflake) ([]Snowflake, bool) {
	membership, exists := s.GetUserMembership(userID)
	if exists {
		return membership.RolesIDs(), true
	}

	return nil, false
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

func (s *Guild) GetUserMembership(userID Snowflake) (*GuildMember, bool) {
	for _, membership := range s.internal.Members {
		if membership.User().ID() == userID {
			return membership, true
		}
	}

	return nil, false
}

func (s *Session) KickUser(guildID, userID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMember(guildID, userID), nil)
}
