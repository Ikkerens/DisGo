package disgo

func onGuildMemberAdd(_ *Session, e GuildMemberAddEvent) {
	objects.guildLock.RLock()
	guild, exists := objects.guilds[e.GuildID]
	objects.guildLock.RUnlock()

	if exists {
		guild.lock.Lock()
		defer guild.lock.Unlock()

		for i, member := range guild.internal.Members {
			if member.User().ID() == e.GuildMember.User().ID() {
				guild.internal.Members[i] = e.GuildMember
				return
			}
		}

		guild.internal.Members = append(guild.internal.Members, e.GuildMember)
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
