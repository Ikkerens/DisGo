package disgo

import "encoding/json"

type snowflake string

type internalUser struct {
	ID            snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Avatar        string    `json:"avatar"`
	Bot           bool      `json:"bot"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	Verified      bool      `json:"verified"`
	EMail         string    `json:"e_mail"`
}

type User struct {
	discordObject internalUser
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.discordObject)
}

func (u *User) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &u.discordObject)
}

func (u *User) Username() string {
	return u.discordObject.Username
}
