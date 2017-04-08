package disgo

import (
	"encoding/json"
	"time"
)

//go:generate go run generate/apimodel/main.go

type Snowflake uint64

type internalUser struct {
	ID            Snowflake `json:"id,string"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	AvatarHash    string    `json:"avatar"`
	Bot           bool      `json:"bot"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	Verified      bool      `json:"verified,omitempty"`
	EMail         string    `json:"e_mail,omitempty"`
}

type internalChannel struct {
	ID                   Snowflake   `json:"id,string"`
	GuildID              Snowflake   `json:"guild_id,string"`
	Name                 string      `json:"name"`
	Type                 string      `json:"type"`
	Position             int         `json:"position"`
	IsPrivate            bool        `json:"is_private"`
	PermissionOverwrites []Overwrite `json:"permission_overwrites"`
	Topic                string      `json:"topic,omitempty"`
	LastMessageID        Snowflake   `json:"last_message_id,string,omitempty"`
	Bitrate              int         `json:"bitrate,omitempty"`
	UserLimit            int         `json:"user_limit,omitempty"`
}

type internalDMChannel struct {
	ID            Snowflake `json:"id,string"`
	IsPrivate     bool      `json:"is_private"`
	Recipient     *User     `json:"recipient"`
	LastMessageID Snowflake `json:"last_message_id,string"`
}

type internalMessage struct {
	ID              Snowflake         `json:"id,string"`
	ChannelID       Snowflake         `json:"channel_id,string"`
	Author          *User             `json:"author"`
	Content         string            `json:"content"`
	Timestamp       time.Time         `json:"timestamp,string"`
	EditedTimestamp time.Time         `json:"edited_timestamp,string,omitempty"`
	TTS             bool              `json:"tts"`
	MentionEveryone bool              `json:"mention_everyone"`
	MentionRoles    []json.RawMessage `json:"mention_roles"`
	Attachments     []Attachment      `json:"attachments"`
	Embeds          []Embed           `json:"embeds"`
	Reactions       []Reaction        `json:"reactions"`
	NOnce           Snowflake         `json:"nonce,string,omitempty"`
	Pinned          bool              `json:"pinned"`
	WebhookID       string            `json:"webhook_id"`
}

type internalReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

type internalEmoji struct {
	ID   Snowflake `json:"id,string,omitempty"`
	Name string    `json:"name"`
}

type internalOverwrite struct {
	ID    Snowflake `json:"id,string"`
	Type  string    `json:"type"`
	Allow int       `json:"allow"`
	Deny  int       `json:"deny"`
}

type internalAttachment struct {
	ID       Snowflake `json:"id,string"`
	Filename string    `json:"filename"`
	Size     int       `json:"size"`
	URL      string    `json:"url"`
	ProxyURL string    `json:"proxy_url"`
	Height   int       `json:"height,omitempty"`
	Width    int       `json:"width,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Type        string         `json:"type,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Timestamp   time.Time      `json:"timestamp,string,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL    string `json:"url"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
