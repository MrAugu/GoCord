package gocord

// Channel - Represents a channel of a guild.
type Channel struct {
	ID                   string                `json:"id"`
	Type                 string                `json:"type"`
	GuildID              string                `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"this_is_for_style_purposes_only_:p"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	Nsfw                 bool                  `json:"nsfw"`
	LastMessageID        string                `json:"last_message_id"`
	Bitrate              int                   `json:"bitrate"`
	UserLimit            int                   `json:"user_limit"`
	Cooldown             int                   `json:"rate_limit_per_user"`
	Recipients           []User                `json:"recipients"`
	Icon                 string                `json:"icon"`
	OwnerID              string                `json:"owner_id"`
	ApplicationID        string                `json:"application_id"`
	ParentID             string                `json:"parent_id"`
	Client               *Client
}

// PermissionOverwrite holds the data of a channel's permission overweries.
type PermissionOverwrite struct {
	ID    string     `json:"id"`
	Type  string     `json:"type"`
	Allow Permission `json:"allow"`
	Deny  Permission `json:"deny"`
	Role  Role
	User  User
}

// 9 Permission Overwrites
