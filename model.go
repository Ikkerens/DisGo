package disgo

import (
	"encoding/json"
	"strconv"
	"time"
)

//go:generate go run generate/apimodel/main.go
//go:generate go run generate/state/main.go

/******************/
/* Resources/Meta */
/******************/

type Snowflake uint64

func ParseSnowflake(str string) (Snowflake, error) {
	intVal, err := strconv.ParseUint(str, 10, 64)
	return Snowflake(intVal), err
}

func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(s), 10))
}

func (s *Snowflake) UnmarshalJSON(b []byte) error {
	var (
		tmp    string
		result uint64
		err    error
	)
	err = json.Unmarshal(b, &tmp)

	if tmp == "" {
		tmp = "0"
	}

	if err == nil {
		result, err = strconv.ParseUint(tmp, 10, 64)
	}

	if err == nil {
		*s = Snowflake(result)
	}

	return err
}

type UnixTimeStamp struct {
	*time.Time
}

func (s UnixTimeStamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Unix())
}

func (s *UnixTimeStamp) UnmarshalJSON(b []byte) error {
	var tmp int64
	err := json.Unmarshal(b, &tmp)

	if err == nil {
		tim := time.Unix(tmp, 0)
		*s = UnixTimeStamp{&tim}
	}

	return err
}

type identifiableObject interface {
	ID() Snowflake
}

type IDObject struct {
	Id Snowflake `json:"id"`
}

func (o *IDObject) ID() Snowflake {
	return o.Id
}

/*********************/
/* Resources/Channel */
/*********************/

type internalChannel struct {
	ID                   Snowflake   `json:"id"`
	GuildID              Snowflake   `json:"guild_id"`
	Name                 string      `json:"name"`
	Type                 string      `json:"type"`
	Position             int         `json:"position"`
	IsPrivate            bool        `json:"is_private"`
	PermissionOverwrites []Overwrite `json:"permission_overwrites"`
	Topic                string      `json:"topic"`
	LastMessageID        Snowflake   `json:"last_message_id,omitempty"`
	Bitrate              int         `json:"bitrate"`
	UserLimit            int         `json:"user_limit"`

	// DMChannel
	Recipient *User `json:"recipient"`
}

type internalMessage struct {
	ID              Snowflake    `json:"id"`
	ChannelID       Snowflake    `json:"channel_id"`
	Author          *User        `json:"author"`
	Content         string       `json:"content"`
	Timestamp       time.Time    `json:"timestamp,string,omitempty"`
	EditedTimestamp time.Time    `json:"edited_timestamp,string,omitempty"`
	TTS             bool         `json:"tts"`
	MentionEveryone bool         `json:"mention_everyone"`
	Mentions        []*User      `json:"mentions"`
	MentionRoles    []*Role      `json:"mention_roles"`
	Attachments     []Attachment `json:"attachments"`
	Embeds          []Embed      `json:"embeds"`
	Reactions       []Reaction   `json:"reactions"`
	NOnce           Snowflake    `json:"nonce"`
	Pinned          bool         `json:"pinned"`
	WebhookID       string       `json:"webhook_id"`
}

type internalReaction struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

type internalOverwrite struct {
	ID    Snowflake `json:"id"`
	Type  string    `json:"type"`
	Allow int       `json:"allow"`
	Deny  int       `json:"deny"`
}

type internalAttachment struct {
	ID       Snowflake `json:"id"`
	Filename string    `json:"filename"`
	Size     int       `json:"size"`
	URL      string    `json:"url"`
	ProxyURL string    `json:"proxy_url"`
	Height   int       `json:"height"`
	Width    int       `json:"width"`
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

/*******************/
/* Resources/Guild */
/*******************/

type internalGuild struct {
	ID                          Snowflake         `json:"id"`
	Name                        string            `json:"name"`
	IconHash                    string            `json:"icon"`
	SplashHash                  string            `json:"splash"`
	OwnerID                     Snowflake         `json:"owner_id"`
	Region                      string            `json:"region"`
	AFKChannelID                Snowflake         `json:"afk_channel_id"`
	AFKTimeout                  int               `json:"afk_timeout"`
	EmbedEnabled                bool              `json:"embed_enabled"`
	EmbedChannelID              Snowflake         `json:"embed_channel_id"`
	VerificationLevel           int               `json:"verification_level"`
	DefaultMessageNotifications int               `json:"default_message_notifications"`
	Roles                       []*Role           `json:"roles"`
	Emojis                      []Emoji           `json:"emojis"`
	Features                    []json.RawMessage `json:"features"`
	MFALevel                    int               `json:"mfa_level"`
	JoinedAt                    time.Time         `json:"joined_at,string"`
	Large                       bool              `json:"large"`
	Unavailable                 bool              `json:"unavailable"`
	MemberCount                 int               `json:"member_count"`
	VoiceStates                 []json.RawMessage `json:"voice_states"`
	Members                     []*GuildMember    `json:"members"`
	Channels                    []*Channel        `json:"channels"`
	Presences                   []Presence        `json:"presences"`
}

type internalGuildMember struct {
	User     *User       `json:"user"`
	Nick     string      `json:"nick,omitempty"`
	RolesIDs []Snowflake `json:"roles"`
	JoinedAt time.Time   `json:"joined_at"`
	Deaf     bool        `json:"deaf"`
	Mute     bool        `json:"mute"`
}

type internalEmoji struct {
	ID            Snowflake `json:"id,omitempty"`
	Name          string    `json:"name"`
	Roles         []*Role   `json:"roles"`
	RequireColons bool      `json:"require_colons"`
	Managed       bool      `json:"managed"`
}

/******************/
/* Resources/User */
/******************/

type internalUser struct {
	ID            Snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	AvatarHash    string    `json:"avatar"`
	Bot           bool      `json:"bot"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	Verified      bool      `json:"verified,omitempty"`
	EMail         string    `json:"e_mail,omitempty"`
}

type internalRole struct {
	ID          Snowflake `json:"id,omitempty"`
	Name        string    `json:"name"`
	Color       int       `json:"color"`
	Hoist       bool      `json:"hoist"`
	Position    int       `json:"position"`
	Permissions int       `json:"permissions"`
	Managed     bool      `json:"managed"`
	Mentionable bool      `json:"mentionable"`
}

type internalPresence struct {
	User    *User       `json:"user"`
	Roles   []Snowflake `json:"roles"`
	Game    Game        `json:"game,omitempty"`
	GuildID Snowflake   `json:"guild_id"`
	Status  string      `json:"status"`
}

type internalGame struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	URL  string `json:"url,omitempty"`
}
