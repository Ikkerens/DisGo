package disgo

import "encoding/json"

type createMessage struct {
	Content string           `json:"content"`
	NOnce   Snowflake        `json:"nonce,omitempty"`
	TTS     bool             `json:"tts,omitempty"`
	File    *json.RawMessage `json:"file,omitempty"`
	Embed   Embed            `json:"embed,omitempty"`
}

func (s *Session) SendMessage(channelID Snowflake, content string) error {
	return s.doHttpPost(EndPointMessages(channelID), createMessage{Content: content})
}

func (s *Session) SendEmbed(channelID Snowflake, embed Embed) error {
	return s.doHttpPost(EndPointMessages(channelID), createMessage{Content: "", Embed: embed})
}

func (s *Session) DeleteMessage(channelID, messageID Snowflake) error {
	return s.doHttpDelete(EndPointMessage(channelID, messageID))
}

func (s *Message) Delete() error {
	return s.session.DeleteMessage(s.discordObject.ChannelID, s.discordObject.ID)
}
