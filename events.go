package disgo

import (
	"encoding/json"
)

type Event interface {
	EventName() string
}

type ReadyEvent struct {
	GatewayVersion  int          `json:"v"`
	User            *User        `json:"user"`
	PrivateChannels []*DMChannel `json:"private_channels"`
	Guilds          []*Guild     `json:"guilds"`
	SessionID       string       `json:"session_id"`
	Servers         []string     `json:"_trace"`
}

func (ReadyEvent) EventName() string {
	return "READY"
}

type ResumedEvent struct {
	Servers []string `json:"_trace"`
}

func (ResumedEvent) EventName() string {
	return "RESUMED"
}

type GuildCreateEvent struct {
	*Guild
}

func (e *GuildCreateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Guild)
}

func (e *GuildCreateEvent) UnmarshalJSON(b []byte) error {
	e.Guild = &Guild{}
	return json.Unmarshal(b, &e.Guild)
}

func (GuildCreateEvent) EventName() string {
	return "GUILD_CREATE"
}

type MessageCreateEvent struct {
	*Message
}

func (e *MessageCreateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Message)
}

func (e *MessageCreateEvent) UnmarshalJSON(b []byte) error {
	e.Message = &Message{}
	return json.Unmarshal(b, &e.Message)
}

func (MessageCreateEvent) EventName() string {
	return "MESSAGE_CREATE"
}

type MessageDeleteEvent struct {
	ID        Snowflake `json:"id"`
	ChannelID Snowflake `json:"channel_id"`
}

func (MessageDeleteEvent) EventName() string {
	return "MESSAGE_DELETE"
}

type TypingStartEvent struct {
	ChannelID Snowflake     `json:"channel_id"`
	UserID    Snowflake     `json:"user_id"`
	Timestamp UnixTimeStamp `json:"timestamp"`
}

func (TypingStartEvent) EventName() string {
	return "TYPING_START"
}
