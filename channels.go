package disgo

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
)

func (s *Session) BuildChannel(guildID Snowflake, name string) *ChannelBuilder {
	return &ChannelBuilder{
		guildID: guildID,
		session: s,

		Name:                 name,
		Type:                 "text",
		PermissionOverwrites: make([]Overwrite, 0),
	}
}

func (s *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", s.ID())
}

func (s *Channel) Guild() *Guild {
	objects.guildLock.RLock()
	defer objects.guildLock.RUnlock()

	return objects.guilds[s.internal.GuildID]
}

func (s *Session) DeleteChannel(channelID Snowflake) error {
	return s.doHttpDelete(EndPointChannel(channelID), nil)
}

func (s *Channel) Delete() error {
	return s.session.DeleteChannel(s.ID())
}

type MessagePrototype struct {
	Content string `json:"content"`
	TTS     bool   `json:"tts"`
	Embed   *Embed `json:"embed,omitempty"`

	FileName string
	File     io.Reader
}

func (s *Session) SendMessage(channelID Snowflake, content string) (*Message, error) {
	return s.SendMessageP(channelID, MessagePrototype{Content: content})
}

func (s *Session) SendEmbed(channelID Snowflake, embed *Embed) (*Message, error) {
	return s.SendMessageP(channelID, MessagePrototype{Embed: embed})
}

func (s *Session) SendMessageP(channelID Snowflake, prototype MessagePrototype) (*Message, error) {
	message := &Message{}

	var err error
	if prototype.File == nil {
		err = s.doHttpPost(EndPointMessages(channelID), &prototype, message)
	} else {
		if prototype.FileName == "" {
			panic("A File was passed to a message without a FileName.")
		}

		var jsonPayload []byte
		jsonPayload, err = json.Marshal(&prototype)
		err = s.doHttMultipartPost(EndPointMessages(channelID), func(writer *multipart.Writer) error {
			writer.WriteField("payload_json", string(jsonPayload))

			if fileW, err := writer.CreateFormFile("file", prototype.FileName); err == nil {
				io.Copy(fileW, prototype.File)
				return nil
			} else {
				return err
			}
		}, message)
	}

	if err != nil {
		return nil, err
	}
	message = objects.registerMessage(message)
	if message.session == nil {
		message.session = s
	}
	return message, nil
}

type editMessageBody struct {
	Content *string `json:"content,omitempty"`
	Embed   *Embed  `json:"embed,omitempty"`
}

func (s *Session) EditMessage(channelID, messageID Snowflake, content string) (*Message, error) {
	return s.editMessageInternal(s.doHttpPatch, EndPointMessage(channelID, messageID), &editMessageBody{Content: &content})
}

func (s *Session) EditEmbed(channelID, messageID Snowflake, embed Embed) (*Message, error) {
	return s.editMessageInternal(s.doHttpPatch, EndPointMessage(channelID, messageID), &editMessageBody{Embed: &embed})
}

func (s *Session) EditEmbeddedMessage(channelID, messageID Snowflake, content string, embed Embed) (*Message, error) {
	return s.editMessageInternal(s.doHttpPatch, EndPointMessage(channelID, messageID), &editMessageBody{Content: &content, Embed: &embed})
}

func (s *Session) editMessageInternal(method func(endPoint EndPoint, body, target interface{}) error, endpoint EndPoint, body *editMessageBody) (*Message, error) {
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

func (s *Session) GetMessage(channelID, messageID Snowflake) (*Message, error) {
	msg, exists := objects.messages[messageID]

	if !exists {
		msg = objects.registerMessage(&IDObject{messageID})
		err := s.doHttpGet(EndPointMessage(channelID, messageID), msg)

		if err != nil {
			objects.messageLock.Lock()
			delete(objects.messages, messageID)
			objects.messageLock.Unlock()
			return nil, err
		}

	}

	return msg, nil
}

type GetMessagesMode int

const (
	GetLastMessages GetMessagesMode = iota
	GetMessagesAround
	GetMessagesBefore
	GetMessagesAfter
)

func (s *Session) GetLastMessages(channelID Snowflake, limit int) ([]*Message, error) {
	return s.GetMessages(channelID, GetLastMessages, 0, limit)
}

func (s *Session) GetMessages(channelID Snowflake, mode GetMessagesMode, target Snowflake, limit int) ([]*Message, error) {
	endPoint := EndPointMessages(channelID)
	limit = int(math.Max(2, math.Min(float64(limit), 100)))

	switch mode {
	case GetLastMessages:
		endPoint.Url += "?"
	case GetMessagesAround:
		endPoint.Url += "?around=" + target.String()
	case GetMessagesBefore:
		endPoint.Url += "?before=" + target.String()
	case GetMessagesAfter:
		endPoint.Url += "?after=" + target.String()
	default:
		panic("Invalid mode parameter passed to Session#GetMessages")
	}

	endPoint.Url += fmt.Sprintf("&limit=%d", limit)
	messages := make([]*Message, 0, limit)

	err := s.doHttpGet(endPoint, &messages)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		objects.registerMessage(message)
		if message.session == nil {
			message.session = s
		}
	}

	return messages, nil
}

func (s *Message) Channel() *Channel {
	objects.channelLock.RLock()
	defer objects.channelLock.RUnlock()

	return objects.channels[s.internal.ChannelID]
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
	endPoint.bucket += "DELETE" // Deleting messages works with a separate rate limit to allow better moderation
	return s.doHttpDelete(endPoint, nil)
}

func (s *Message) Delete() error {
	return s.session.DeleteMessage(s.internal.ChannelID, s.internal.ID)
}

func (s *Session) BulkDeleteMessages(channelID Snowflake, ids []Snowflake) error {
	return s.doHttpPost(EndPointMessageBulkDelete(channelID), struct {
		Messages []Snowflake `json:"messages"`
	}{ids}, nil)
}

func (s *Session) PinMessage(channelID, messageID Snowflake) error {
	return s.doHttpPut(EndPointChannelPin(channelID, messageID), nil)
}

func (s *Message) Pin() error {
	return s.session.PinMessage(s.internal.ChannelID, s.internal.ID)
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
