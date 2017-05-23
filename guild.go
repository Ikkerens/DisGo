package disgo

func (s *Session) AddGuildMemberRole(guildID, userID, roleID Snowflake) error {
	return s.doHttpPut(EndPointGuildMemberRoles(guildID, userID, roleID), nil)
}

func (s *Session) RemoveGuildMemberRole(guildID, userID, roleID Snowflake) error {
	return s.doHttpDelete(EndPointGuildMemberRoles(guildID, userID, roleID), nil)
}
