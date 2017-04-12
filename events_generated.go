package disgo

// Warning: This file has been automatically generated by generate/eventmethods/main.go
// Do NOT make changes here, instead adapt events.go and run go generate

import (
	"encoding/json"
	"github.com/slf4go/logger"
)

func allocateEvent(eventName string) *Event {
	var event Event

	switch eventName {
	case "GUILD_CREATE":
		event = &GuildCreateEvent{}
	case "MESSAGE_CREATE":
		event = &MessageCreateEvent{}
	case "MESSAGE_DELETE":
		event = &MessageDeleteEvent{}
	case "PRESENCE_UPDATE":
		event = &PresenceUpdateEvent{}
	case "READY":
		event = &ReadyEvent{}
	case "RESUMED":
		event = &ResumedEvent{}
	case "TYPING_START":
		event = &TypingStartEvent{}
	default:
		logger.Errorf("Event with name '%s' was dispatched by Discord, but we don't know this event. (DisGo outdated?)", eventName)
		return nil
	}

	return &event
}

func (*GuildCreateEvent) eventName() string {
	return "GUILD_CREATE"
}

func (e *GuildCreateEvent) setSession(s *Session) {
	e.Guild = s.registerGuild(e.Guild)
}

// MarshalJSON is used to make sure the embedded object of this event is Marshalled, not the event itself
func (e *GuildCreateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Guild)
}

// UnmarshalJSON is used to make sure the embedded object of this event is Unmarshalled, not the event itself
func (e *GuildCreateEvent) UnmarshalJSON(b []byte) error {
	e.Guild = &Guild{}
	return json.Unmarshal(b, &e.Guild)
}

func (*MessageCreateEvent) eventName() string {
	return "MESSAGE_CREATE"
}

func (e *MessageCreateEvent) setSession(s *Session) {
	e.Message = s.registerMessage(e.Message)
}

// MarshalJSON is used to make sure the embedded object of this event is Marshalled, not the event itself
func (e *MessageCreateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Message)
}

// UnmarshalJSON is used to make sure the embedded object of this event is Unmarshalled, not the event itself
func (e *MessageCreateEvent) UnmarshalJSON(b []byte) error {
	e.Message = &Message{}
	return json.Unmarshal(b, &e.Message)
}

func (*MessageDeleteEvent) eventName() string {
	return "MESSAGE_DELETE"
}

func (e *MessageDeleteEvent) setSession(s *Session) {
}

func (*PresenceUpdateEvent) eventName() string {
	return "PRESENCE_UPDATE"
}

func (e *PresenceUpdateEvent) setSession(s *Session) {
}

// MarshalJSON is used to make sure the embedded object of this event is Marshalled, not the event itself
func (e *PresenceUpdateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Presence)
}

// UnmarshalJSON is used to make sure the embedded object of this event is Unmarshalled, not the event itself
func (e *PresenceUpdateEvent) UnmarshalJSON(b []byte) error {
	e.Presence = &Presence{}
	return json.Unmarshal(b, &e.Presence)
}

func (*ReadyEvent) eventName() string {
	return "READY"
}

func (e *ReadyEvent) setSession(s *Session) {
	e.User = s.registerUser(e.User)
	for i, p := range e.Guilds {
		e.Guilds[i] = s.registerGuild(p)
	}
}

func (*ResumedEvent) eventName() string {
	return "RESUMED"
}

func (e *ResumedEvent) setSession(s *Session) {
}

func (*TypingStartEvent) eventName() string {
	return "TYPING_START"
}

func (e *TypingStartEvent) setSession(s *Session) {
}
