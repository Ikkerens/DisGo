package disgo

import "encoding/json"

type createMessage struct {
	Content string           `json:"content"`
	NOnce   Snowflake        `json:"nonce,omitempty"`
	TTS     bool             `json:"tts,omitempty"`
	File    *json.RawMessage `json:"file,omitempty"`
	Embed   *Embed           `json:"embed,omitempty"`
}

func (s *Session) SendMessage(channelID Snowflake, content string) (message *Message, err error) {
	message = &Message{}
	err = s.doHttpPost(EndPointMessages(channelID), createMessage{Content: content}, message)
	s.registerMessage(message)
	return
}

func (s *Session) SendEmbed(channelID Snowflake, embed Embed) (message *Message, err error) {
	message = &Message{}
	err = s.doHttpPost(EndPointMessages(channelID), createMessage{Content: "", Embed: &embed}, message)
	s.registerMessage(message)
	return
}

func (s *Session) DeleteMessage(channelID, messageID Snowflake) error {
	return s.doHttpDelete(EndPointMessage(channelID, messageID), nil)
}

func (s *Message) Delete() error {
	return s.session.DeleteMessage(s.internal.ChannelID, s.internal.ID)
}
