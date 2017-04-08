package disgo

import "encoding/json"

type Event interface {
	EventName() string
}

type ReadyEvent struct {
	GatewayVersion  int               `json:"v"`
	User            *User             `json:"user"`
	PrivateChannels []*DMChannel      `json:"private_channels"`
	Guilds          []json.RawMessage `json:"guilds"`
	SessionID       string            `json:"session_id"`
	Servers         []string          `json:"_trace"`
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
