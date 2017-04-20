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

type PresenceUpdateEvent struct {
	*Presence
}

type GuildCreateEvent struct {
	*Guild
}

type MessageCreateEvent struct {
	*Message
}

type MessageDeleteEvent struct {
	ID        Snowflake `json:"id"`
	ChannelID Snowflake `json:"channel_id"`
}

type TypingStartEvent struct {
	ChannelID Snowflake     `json:"channel_id"`
	UserID    Snowflake     `json:"user_id"`
	Timestamp UnixTimeStamp `json:"timestamp"`
}
