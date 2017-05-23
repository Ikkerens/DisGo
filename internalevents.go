package disgo

func onGuildMemberUpdate(_ *Session, e GuildMemberUpdateEvent) {
	guild, exists := objects.guilds[e.GuildID]
	if exists {
		members := guild.Members()
		for i := range members {
			if members[i].User().ID() == e.User.ID() {
				internal := members[i].internal
				internal.RolesIDs = e.Roles
				internal.Nick = e.Nick
			}
		}
	}
}
