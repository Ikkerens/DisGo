package disgo

//go:generate go run generate/apimodel/main.go

type Snowflake uint64

type internalUser struct {
	ID            Snowflake `json:"id,string"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Avatar        string    `json:"avatar"`
	Bot           bool      `json:"bot"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	Verified      bool      `json:"verified,omitempty"`
	EMail         string    `json:"e_mail,omitempty"`
}

type internalDMChannel struct {
	ID            Snowflake `json:"id"`
	IsPrivate     bool      `json:"is_private"`
	Recipient     *User     `json:"recipient"`
	LastMessageID Snowflake `json:"last_message_id"`
}