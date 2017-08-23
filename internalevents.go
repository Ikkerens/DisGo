package disgo

func registerInternalEvents(session *Session) {
	session.registerEventHandler(onReady, false)
	session.registerEventHandler(onGuildCreate, false)
	session.registerEventHandler(onChannelCreate, false)
	session.registerEventHandler(onChannelDelete, false)
	session.registerEventHandler(onGuildMemberUpdate, false)
	session.registerEventHandler(onGuildMemberAdd, false)
	session.registerEventHandler(onGuildMemberRemove, false)
}

func onReady(_ *Session, e ReadyEvent) {
	for _, guild := range e.Guilds {
		for _, channel := range guild.Channels() {
			channel.internal.GuildID = guild.internal.ID
		}
	}
}

func onGuildCreate(_ *Session, e GuildCreateEvent) {
	for _, channel := range e.Channels() {
		channel.internal.GuildID = e.internal.ID
	}
}

func onGuildMemberAdd(_ *Session, e GuildMemberAddEvent) {
	objects.guildLock.RLock()
	guild, exists := objects.guilds[e.GuildID]
	objects.guildLock.RUnlock()

	if exists {
		guild.lock.Lock()
		defer guild.lock.Unlock()

		for i, member := range guild.internal.Members {
			if member.internal.User.internal.ID == e.GuildMember.internal.User.internal.ID {
				guild.internal.Members[i] = e.GuildMember
				return
			}
		}

		guild.internal.Members = append(guild.internal.Members, e.GuildMember)
	}
}

func onGuildMemberRemove(_ *Session, e GuildMemberRemoveEvent) {
	objects.guildLock.RLock()
	guild, exists := objects.guilds[e.GuildID]
	objects.guildLock.RUnlock()

	if exists {
		guild.lock.Lock()
		defer guild.lock.Unlock()

		for i, member := range guild.internal.Members {
			if member.internal.User.internal.ID == e.User.internal.ID {
				guild.internal.Members = append(guild.internal.Members[:i], guild.internal.Members[i+1:]...)
				return
			}
		}
	}
}

func onGuildMemberUpdate(_ *Session, e GuildMemberUpdateEvent) {
	objects.guildLock.RLock()
	guild, exists := objects.guilds[e.GuildID]
	objects.guildLock.RUnlock()

	if exists {
		guild.lock.Lock()
		defer guild.lock.Unlock()

		membership, exists := guild.GetUserMembership(e.User.ID())

		if exists {
			membership.internal.RolesIDs = e.Roles
			membership.internal.Nick = e.Nick
		}
	}
}

func onChannelCreate(_ *Session, e ChannelCreateEvent) {
	objects.guildLock.RLock()
	guild, exists := objects.guilds[e.GuildID()]
	objects.guildLock.RUnlock()

	if exists {
		guild.lock.Lock()
		defer guild.lock.Unlock()

		guild.internal.Channels = append(guild.internal.Channels, e.Channel)
	}
}

func onChannelDelete(_ *Session, e ChannelDeleteEvent) {
	objects.guildLock.RLock()
	guild, exists := objects.guilds[e.GuildID()]
	objects.guildLock.RUnlock()

	if exists {
		guild.lock.Lock()
		defer guild.lock.Unlock()

		index := -1
		for i, channel := range guild.internal.Channels {
			if channel.ID() == e.ID() {
				index = i
				break
			}
		}

		if index != -1 {
			guild.internal.Channels = append(guild.internal.Channels[:index], guild.internal.Channels[index+1:]...)
		}
	}
}
