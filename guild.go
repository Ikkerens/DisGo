package disgo

func (s *Guild) BuildChannel(name string) *ChannelBuilder {
	return s.session.BuildChannel(s.internal.ID, name)
}

func (s *Session) AddGuildMemberRole(guildID, userID, roleID Snowflake) error {
	return s.doHttpPut(EndPointGuildMemberRoles(guildID, userID, roleID), nil)
}

func (s *Session) RemoveGuildMemberRole(guildID, userID, roleID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMemberRoles(guildID, userID, roleID), nil)
}

func (s *Guild) GetUserMembership(userID Snowflake) (*GuildMember, bool) {
	for _, membership := range s.internal.Members {
		if membership.User().ID() == userID {
			return membership, true
		}
	}

	return nil, false
}

func (s *Guild) GetUserRoles(userID Snowflake) ([]Snowflake, bool) {
	membership, exists := s.GetUserMembership(userID)
	if exists {
		return membership.RolesIDs(), true
	}

	return nil, false
}

func (s *Session) KickUser(guildID, userID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMember(guildID, userID), nil)
}
