package disgo

//go:generate go run generate/eventmethods/main.go

type ReadyEvent struct {
	GatewayVersion int      `json:"v"`
	User           *User    `json:"user"`
	Guilds         []*Guild `json:"guilds"`
	SessionID      string   `json:"session_id"`
	Servers        []string `json:"_trace"`
}

type ResumedEvent struct {
	Servers []string `json:"_trace"`
}

type ChannelCreateEvent struct {
	*Channel
}

type ChannelUpdateEvent struct {
	*Channel
}

type ChannelDeleteEvent struct {
	*Channel
}

type GuildCreateEvent struct {
	*Guild
}

type GuildUpdateEvent struct {
	*Guild
}

type GuildDeleteEvent struct {
	*Guild
}

type GuildBanAddEvent struct {
	*User
	GuildID Snowflake `json:"guild_id"`
}

type GuildBanRemoveEvent struct {
	*User
	GuildID Snowflake `json:"guild_id"`
}

type GuildEmojisUpdateEvent struct {
	GuildID Snowflake `json:"guild_id"`
	Emojis  []Emoji   `json:"emojis"`
}

type GuildIntegrationsUpdateEvent struct {
	GuildID Snowflake `json:"guild_id"`
}

type GuildMemberAddEvent struct {
	*GuildMember
	GuildID Snowflake `json:"guild_id"`
}

type GuildMemberRemoveEvent struct {
	GuildID Snowflake `json:"guild_id"`
	User    *User     `json:"user"`
}

type GuildMemberUpdateEvent struct {
	GuildID Snowflake   `json:"guild_id"`
	Roles   []Snowflake `json:"roles"`
	User    *User       `json:"user"`
	Nick    string      `json:"nick"`
}

type GuildMembersChunkEvent struct {
	GuildID Snowflake     `json:"guild_id"`
	Members []GuildMember `json:"members"`
}

type GuildRoleCreateEvent struct {
	GuildID Snowflake `json:"guild_id"`
	Role    *Role     `json:"role"`
}

type GuildRoleUpdateEvent struct {
	GuildID Snowflake `json:"guild_id"`
	Role    *Role     `json:"role"`
}

type GuildRoleDeleteEvent struct {
	GuildID Snowflake `json:"guild_id"`
	RoleID  Snowflake `json:"role_id"`
}

type MessageCreateEvent struct {
	*Message
}

type MessageUpdateEvent struct {
	*Message
}

type MessageDeleteEvent struct {
	ID        Snowflake `json:"id"`
	ChannelID Snowflake `json:"channel_id"`
}

type MessageDeleteBulkEvent struct {
	IDs       []Snowflake `json:"ids"`
	ChannelID Snowflake   `json:"channel_id"`
}

type PresenceUpdateEvent struct {
	*Presence
}

type TypingStartEvent struct {
	ChannelID Snowflake     `json:"channel_id"`
	UserID    Snowflake     `json:"user_id"`
	Timestamp UnixTimeStamp `json:"timestamp"`
}

type UserUpdateEvent struct {
	*User
}
