package disgo

func (s *Session) registerGuild(guild *Guild) *Guild {
	if snowflake, exists := s.objects[guild.ID()]; exists {
		if registered, isGuild := snowflake.(*Guild); isGuild {
			// TODO Merge data
			return registered
		} else {
			panic("Discord sent us a duplicate snowflake with different types")
		}
	} else {
		s.objects[guild.ID()] = guild
		guild.session = s
		return guild
	}
}

func (s *Session) registerUser(user *User) *User {
	if snowflake, exists := s.objects[user.ID()]; exists {
		if registered, isUser := snowflake.(*User); isUser {
			// TODO Merge data
			return registered
		} else {
			panic("Discord sent us a duplicate snowflake with different types")
		}
	} else {
		s.objects[user.ID()] = user
		user.session = s
		return user
	}
}

func (s *Session) registerChannel(channel *Channel) *Channel {
	if snowflake, exists := s.objects[channel.ID()]; exists {
		if registered, isChannel := snowflake.(*Channel); isChannel {
			// TODO Merge data
			return registered
		} else {
			panic("Discord sent us a duplicate snowflake with different types")
		}
	} else {
		s.objects[channel.ID()] = channel
		channel.session = s
		return channel
	}
}

func (s *Session) registerMessage(message *Message) *Message {
	if snowflake, exists := s.objects[message.ID()]; exists {
		if registered, isMessage := snowflake.(*Message); isMessage {
			// TODO Merge data
			return registered
		} else {
			panic("Discord sent us a duplicate snowflake with different types")
		}
	} else {
		s.objects[message.ID()] = message
		message.session = s
		return message
	}
}
