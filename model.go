package disgo

import "encoding/json"

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

type User struct {
	discordObject *internalUser
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.discordObject)
}

func (u *User) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &u.discordObject)
}

func (u *User) ID() Snowflake {
	return u.discordObject.ID
}

func (u *User) Username() string {
	return u.discordObject.Username
}
