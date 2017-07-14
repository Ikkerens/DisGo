package disgo

func (e *MessageCreateEvent) Reply(content string) (*Message, error) {
	return e.session.SendMessage(e.ChannelID(), content)
}

func (e *MessageCreateEvent) Channel() *Channel {
	objects.channelLock.RLock()
	defer objects.channelLock.RUnlock()

	return objects.channels[e.internal.ChannelID]
}
