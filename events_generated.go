package disgo

// Warning: This file has been automatically generated by generate/eventmethods/main.go
// Do NOT make changes here, instead adapt events.go and run go generate

import (
	"encoding/json"
	"github.com/slf4go/logger"
)

func allocateEvent(eventName string) *Event {
	var event Event

	// Because encoding/json doesn't initialise embbeded struct pointers properly, we'll also initialise them here
	switch eventName {
	case "CHANNEL_CREATE":
		event = &ChannelCreateEvent{Channel: &Channel{}}
	case "CHANNEL_DELETE":
		event = &ChannelDeleteEvent{Channel: &Channel{}}
	case "CHANNEL_UPDATE":
		event = &ChannelUpdateEvent{Channel: &Channel{}}
	case "GUILD_BAN_ADD":
		event = &GuildBanAddEvent{User: &User{}}
	case "GUILD_BAN_REMOVE":
		event = &GuildBanRemoveEvent{User: &User{}}
	case "GUILD_CREATE":
		event = &GuildCreateEvent{Guild: &Guild{}}
	case "GUILD_DELETE":
		event = &GuildDeleteEvent{Guild: &Guild{}}
	case "GUILD_EMOJIS_UPDATE":
		event = &GuildEmojisUpdateEvent{}
	case "GUILD_INTEGRATIONS_UPDATE":
		event = &GuildIntegrationsUpdateEvent{}
	case "GUILD_MEMBER_ADD":
		event = &GuildMemberAddEvent{GuildMember: &GuildMember{}}
	case "GUILD_MEMBER_REMOVE":
		event = &GuildMemberRemoveEvent{}
	case "GUILD_MEMBER_UPDATE":
		event = &GuildMemberUpdateEvent{}
	case "GUILD_MEMBERS_CHUNK":
		event = &GuildMembersChunkEvent{}
	case "GUILD_ROLE_CREATE":
		event = &GuildRoleCreateEvent{}
	case "GUILD_ROLE_DELETE":
		event = &GuildRoleDeleteEvent{}
	case "GUILD_ROLE_UPDATE":
		event = &GuildRoleUpdateEvent{}
	case "GUILD_UPDATE":
		event = &GuildUpdateEvent{Guild: &Guild{}}
	case "MESSAGE_CREATE":
		event = &MessageCreateEvent{Message: &Message{}}
	case "MESSAGE_DELETE_BULK":
		event = &MessageDeleteBulkEvent{}
	case "MESSAGE_DELETE":
		event = &MessageDeleteEvent{}
	case "MESSAGE_REACTION_ADD":
		event = &MessageReactionAddEvent{}
	case "MESSAGE_REACTION_REMOVE":
		event = &MessageReactionRemoveEvent{}
	case "MESSAGE_UPDATE":
		event = &MessageUpdateEvent{Message: &Message{}}
	case "PRESENCE_UPDATE":
		event = &PresenceUpdateEvent{Presence: &Presence{}}
	case "READY":
		event = &ReadyEvent{}
	case "RESUMED":
		event = &ResumedEvent{}
	case "TYPING_START":
		event = &TypingStartEvent{}
	case "USER_UPDATE":
		event = &UserUpdateEvent{User: &User{}}
	default:
		logger.Errorf("Event with name '%s' was dispatched by Discord, but we don't know this event. (DisGo outdated?)", eventName)
		return nil
	}

	return &event
}

func (*ChannelCreateEvent) eventName() string {
	return "CHANNEL_CREATE"
}

func (e *ChannelCreateEvent) setSession(s *Session) {
	e.Channel.session = s
}

func (*ChannelDeleteEvent) eventName() string {
	return "CHANNEL_DELETE"
}

func (e *ChannelDeleteEvent) setSession(s *Session) {
	e.Channel.session = s
}

func (*ChannelUpdateEvent) eventName() string {
	return "CHANNEL_UPDATE"
}

func (e *ChannelUpdateEvent) setSession(s *Session) {
	e.Channel.session = s
}

func (*GuildBanAddEvent) eventName() string {
	return "GUILD_BAN_ADD"
}

func (e *GuildBanAddEvent) setSession(s *Session) {
	e.User.session = s
}

func (*GuildBanRemoveEvent) eventName() string {
	return "GUILD_BAN_REMOVE"
}

func (e *GuildBanRemoveEvent) setSession(s *Session) {
	e.User.session = s
}

func (*GuildCreateEvent) eventName() string {
	return "GUILD_CREATE"
}

func (e *GuildCreateEvent) setSession(s *Session) {
	e.Guild.session = s
}

func (*GuildDeleteEvent) eventName() string {
	return "GUILD_DELETE"
}

func (e *GuildDeleteEvent) setSession(s *Session) {
	e.Guild.session = s
}

func (*GuildEmojisUpdateEvent) eventName() string {
	return "GUILD_EMOJIS_UPDATE"
}

func (e *GuildEmojisUpdateEvent) setSession(s *Session) {
}

func (*GuildIntegrationsUpdateEvent) eventName() string {
	return "GUILD_INTEGRATIONS_UPDATE"
}

func (e *GuildIntegrationsUpdateEvent) setSession(s *Session) {
}

func (*GuildMemberAddEvent) eventName() string {
	return "GUILD_MEMBER_ADD"
}

func (e *GuildMemberAddEvent) setSession(s *Session) {
	e.GuildMember.session = s
}

func (*GuildMemberRemoveEvent) eventName() string {
	return "GUILD_MEMBER_REMOVE"
}

func (e *GuildMemberRemoveEvent) setSession(s *Session) {
	e.User.session = s
}

func (*GuildMemberUpdateEvent) eventName() string {
	return "GUILD_MEMBER_UPDATE"
}

func (e *GuildMemberUpdateEvent) setSession(s *Session) {
	e.User.session = s
}

func (*GuildMembersChunkEvent) eventName() string {
	return "GUILD_MEMBERS_CHUNK"
}

func (e *GuildMembersChunkEvent) setSession(s *Session) {
}

func (*GuildRoleCreateEvent) eventName() string {
	return "GUILD_ROLE_CREATE"
}

func (e *GuildRoleCreateEvent) setSession(s *Session) {
	e.Role.session = s
}

func (*GuildRoleDeleteEvent) eventName() string {
	return "GUILD_ROLE_DELETE"
}

func (e *GuildRoleDeleteEvent) setSession(s *Session) {
}

func (*GuildRoleUpdateEvent) eventName() string {
	return "GUILD_ROLE_UPDATE"
}

func (e *GuildRoleUpdateEvent) setSession(s *Session) {
	e.Role.session = s
}

func (*GuildUpdateEvent) eventName() string {
	return "GUILD_UPDATE"
}

func (e *GuildUpdateEvent) setSession(s *Session) {
	e.Guild.session = s
}

func (*MessageCreateEvent) eventName() string {
	return "MESSAGE_CREATE"
}

func (e *MessageCreateEvent) setSession(s *Session) {
	e.Message.session = s
}

func (*MessageDeleteBulkEvent) eventName() string {
	return "MESSAGE_DELETE_BULK"
}

func (e *MessageDeleteBulkEvent) setSession(s *Session) {
}

func (*MessageDeleteEvent) eventName() string {
	return "MESSAGE_DELETE"
}

func (e *MessageDeleteEvent) setSession(s *Session) {
}

func (*MessageReactionAddEvent) eventName() string {
	return "MESSAGE_REACTION_ADD"
}

func (e *MessageReactionAddEvent) setSession(s *Session) {
}

func (*MessageReactionRemoveEvent) eventName() string {
	return "MESSAGE_REACTION_REMOVE"
}

func (e *MessageReactionRemoveEvent) setSession(s *Session) {
}

func (*MessageUpdateEvent) eventName() string {
	return "MESSAGE_UPDATE"
}

func (e *MessageUpdateEvent) setSession(s *Session) {
	e.Message.session = s
}

func (*PresenceUpdateEvent) eventName() string {
	return "PRESENCE_UPDATE"
}

func (e *PresenceUpdateEvent) setSession(s *Session) {
}

func (*ReadyEvent) eventName() string {
	return "READY"
}

func (e *ReadyEvent) setSession(s *Session) {
	e.User.session = s
	for _, item := range e.Guilds {
		item.session = s
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

func (*UserUpdateEvent) eventName() string {
	return "USER_UPDATE"
}

func (e *UserUpdateEvent) setSession(s *Session) {
	e.User.session = s
}
