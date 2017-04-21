package disgo

func (e *MessageCreateEvent) Reply(content string) (*Message, error) {
	return e.session.SendMessage(e.ChannelID(), content)
}
