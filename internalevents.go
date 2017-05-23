package disgo

func onGuildMemberAdd(_ *Session, e GuildMemberAddEvent) {
	guild, exists := objects.guilds[e.GuildID]
	if exists {
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
	guild, exists := objects.guilds[e.GuildID]
	if exists {
		for i := range guild.internal.Members {
			if guild.internal.Members[i].User().ID() == e.User.ID() {
				internal := guild.internal.Members[i].internal
				internal.RolesIDs = e.Roles
				internal.Nick = e.Nick
			}
		}
	}
}
