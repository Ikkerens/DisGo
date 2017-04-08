package disgo

import "encoding/json"

type createMessage struct {
	Content string           `json:"content"`
	NOnce   Snowflake        `json:"nonce,omitempty"`
	TTS     bool             `json:"tts,omitempty"`
	File    *json.RawMessage `json:"file,omitempty"`
	Embed   Embed            `json:"embed,omitempty"`
}

func (s *Session) SendMessage(channelID Snowflake, content string) (err error) {
	err = s.doHttpPost(EndPointMessages(channelID), createMessage{Content: content})
	return
}

func (s *Session) SendEmbed(channelID Snowflake, embed Embed) (err error) {
	err = s.doHttpPost(EndPointMessages(channelID), createMessage{Content: "", Embed: embed})
	return
}

func (s *Session) DeleteMessage(channelID, messageID Snowflake) (err error) {
	err = s.doHttpDelete(EndPointMessage(channelID, messageID))
	return
}
