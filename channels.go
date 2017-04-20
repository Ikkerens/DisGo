package disgo

type createMessage struct {
	Content string `json:"content"`
	Embed   *Embed `json:"embed,omitempty"`
}

func (s *Session) SendMessage(channelID Snowflake, content string) (*Message, error) {
	message := &Message{}
	err := s.doHttpPost(EndPointMessages(channelID), createMessage{Content: content}, message)
	if err != nil {
		return nil, err
	}
	return objects.registerMessage(message), nil
}

func (s *Session) SendEmbed(channelID Snowflake, embed Embed) (*Message, error) {
	message := &Message{}
	err := s.doHttpPost(EndPointMessages(channelID), createMessage{Content: "", Embed: &embed}, message)
	if err != nil {
		return nil, err
	}
	return objects.registerMessage(message), nil
}

func (s *Session) DeleteMessage(channelID, messageID Snowflake) error {
	endPoint := EndPointMessage(channelID, messageID)
	endPoint.Bucket += "DELETE" // Deleting messages works with a separate ratelimit to allow better moderation
	return s.doHttpDelete(endPoint, nil)
}

func (s *Message) Delete() error {
	return s.session.DeleteMessage(s.internal.ChannelID, s.internal.ID)
}
