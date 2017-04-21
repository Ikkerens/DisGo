package disgo

import "fmt"

type sendMessageBody struct {
	Content string `json:"content"`
	Embed   *Embed `json:"embed,omitempty"`
}

func (s *Session) SendMessage(channelID Snowflake, content string) (*Message, error) {
	return s.sendMessageInternal(s.doHttpPost, EndPointMessages(channelID), &sendMessageBody{Content: content})
}

func (s *Session) SendEmbed(channelID Snowflake, embed Embed) (*Message, error) {
	return s.sendMessageInternal(s.doHttpPost, EndPointMessages(channelID), &sendMessageBody{Content: "", Embed: &embed})
}

func (s *Session) SendEmbeddedMessage(channelID Snowflake, content string, embed Embed) (*Message, error) {
	return s.sendMessageInternal(s.doHttpPost, EndPointMessages(channelID), &sendMessageBody{Content: content, Embed: &embed})
}

func (s *Session) EditMessage(channelID, messageID Snowflake, content string) (*Message, error) {
	return s.sendMessageInternal(s.doHttpPatch, EndPointMessage(channelID, messageID), &sendMessageBody{Content: content})
}

func (s *Session) EditEmbed(channelID, messageID Snowflake, embed Embed) (*Message, error) {
	return s.sendMessageInternal(s.doHttpPatch, EndPointMessage(channelID, messageID), &sendMessageBody{Embed: &embed})
}

func (s *Session) EditEmbeddedMessage(channelID, messageID Snowflake, content string, embed Embed) (*Message, error) {
	return s.sendMessageInternal(s.doHttpPatch, EndPointMessage(channelID, messageID), &sendMessageBody{Content: content, Embed: &embed})
}

func (s *Session) sendMessageInternal(method func(endPoint EndPoint, body, target interface{}) error, endpoint EndPoint, body *sendMessageBody) (*Message, error) {
	message := &Message{}
	err := method(endpoint, body, message)
	if err != nil {
		return nil, err
	}
	message = objects.registerMessage(message)
	if message.session == nil {
		message.session = s
	}
	return message, nil
}

func (s *Message) Edit(content string) (err error) {
	_, err = s.session.EditMessage(s.internal.ChannelID, s.internal.ID, content)
	return
}

func (s *Message) EditEmbed(embed Embed) (err error) {
	_, err = s.session.EditEmbed(s.internal.ChannelID, s.internal.ID, embed)
	return
}

func (s *Message) EditEmbeddedMessage(content string, embed Embed) (err error) {
	_, err = s.session.EditEmbeddedMessage(s.internal.ChannelID, s.internal.ID, content, embed)
	return
}

func (s *Session) DeleteMessage(channelID, messageID Snowflake) error {
	endPoint := EndPointMessage(channelID, messageID)
	endPoint.bucket += "DELETE" // Deleting messages works with a separate ratelimit to allow better moderation
	return s.doHttpDelete(endPoint, nil)
}

func (s *Message) Delete() error {
	return s.session.DeleteMessage(s.internal.ChannelID, s.internal.ID)
}

func (s *Session) MessageAddReaction(channelID, messageID Snowflake, emoji string) error {
	endPoint := EndPointOwnReaction(channelID, messageID)
	endPoint.Url = fmt.Sprintf(endPoint.Url, emoji)
	endPoint.resetTime = 300
	return s.doHttpPut(endPoint, nil)
}

func (s *Message) AddReaction(emoji string) error {
	return s.session.MessageAddReaction(s.internal.ChannelID, s.internal.ID, emoji)
}

func (s *Session) MessageDeleteOwnReaction(channelID, messageID Snowflake, emoji string) error {
	endPoint := EndPointOwnReaction(channelID, messageID)
	endPoint.Url = fmt.Sprintf(endPoint.Url, emoji)
	endPoint.resetTime = 250
	return s.doHttpDelete(endPoint, nil)
}

func (s *Session) MessageDeleteReaction(channelID, messageID, userID Snowflake, emoji string) error {
	endPoint := EndPointReaction(channelID, messageID, userID)
	endPoint.Url = fmt.Sprintf(endPoint.Url, emoji)
	endPoint.resetTime = 250
	return s.doHttpDelete(endPoint, nil)
}

func (s *Session) MessageDeleteAllReactions(channelID, messageID Snowflake) error {
	return s.doHttpDelete(EndPointReactions(channelID, messageID), nil)
}

func (s *Message) DeleteReaction(userID Snowflake, emoji string) error {
	return s.session.MessageDeleteReaction(s.internal.ChannelID, s.internal.ID, userID, emoji)
}

func (s *Message) DeleteOwnReaction(emoji string) error {
	return s.session.MessageDeleteOwnReaction(s.internal.ChannelID, s.internal.ID, emoji)
}

func (s *Message) DeleteAllReactions() error {
	return s.session.MessageDeleteAllReactions(s.internal.ChannelID, s.internal.ID)
}
