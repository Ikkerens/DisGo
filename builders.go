package disgo

type ChannelBuilder struct {
	session *Session
	guildID Snowflake

	Name                 string      `json:"name"`
	Type                 string      `json:"type,omitempty"`
	Bitrate              int         `json:"bitrate,omitempty"`
	UserLimit            int         `json:"user_limit,omitempty"`
	PermissionOverwrites []Overwrite `json:"permission_overwrites"`
}

func (b *ChannelBuilder) AddMemberOverwrite(id Snowflake, allow, deny int) *ChannelBuilder {
	b.PermissionOverwrites = append(b.PermissionOverwrites, Overwrite{
		ID:    id,
		Type:  "member",
		Allow: allow,
		Deny:  deny,
	})
	return b
}

func (b *ChannelBuilder) AddRoleOverwrite(id Snowflake, allow, deny int) *ChannelBuilder {
	b.PermissionOverwrites = append(b.PermissionOverwrites, Overwrite{
		ID:    id,
		Type:  "role",
		Allow: allow,
		Deny:  deny,
	})
	return b
}

func (b *ChannelBuilder) Create() (*Channel, error) {
	channel := &Channel{}
	err := b.session.doHttpPost(EndPointGuildChannels(b.guildID), b, channel)
	if err != nil {
		return nil, err
	}
	channel = objects.registerChannel(channel)
	if channel.session == nil {
		channel.session = b.session
	}
	return channel, nil
}
