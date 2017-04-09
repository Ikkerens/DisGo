package disgo

func (s *Session) registerGuild(guild *Guild) *Guild {
	snowflake, exists := s.objects[guild.ID()]
	if exists {
		registered, isGuild := snowflake.(*Guild)
		if isGuild {
			// TODO Merge data
			return registered
		} else {
			panic("Discord sent us a duplicate snowflake with different types")
		}
	} else {
		s.objects[guild.ID()] = guild
		return guild
	}
}

func (s *Session) registerUser(user *User) *User {
	snowflake, exists := s.objects[user.ID()]
	if exists {
		registered, isUser := snowflake.(*User)
		if isUser {
			// TODO Merge data
			return registered
		} else {
			panic("Discord sent us a duplicate snowflake with different types")
		}
	} else {
		s.objects[user.ID()] = user
		return user
	}
}
