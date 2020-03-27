package gocord

// Channel - Represents a channel of a guild.
type Channel struct {
	ID                   string                 `json:"id"`
	Type                 string                 `json:"type"`
	GuildID              string                 `json:"guild_id"`
	Position             int                    `json:"position"`
	PermissionOverwrites []PermissionOverwrites `json:"pqm"`
	Name                 string                 `json:"name"`
	Topic                string                 `json:"topic"`
	Nsfw                 bool                   `json:"nsfw"`
	LastMessageID        string                 `json:"last_message_id"`
	Bitrate              int                    `json:"bitrate"`
	UserLimit            int                    `json:"user_limit"`
	Cooldown             int                    `json:"rate_limit_per_user"`
	Recipients           []User                 `json:"recipients"`
	Icon                 string                 `json:"icon"`
	OwnerID              string                 `json:"owner_id"`
	ApplicationID        string                 `json:"application_id"`
	ParentID             string                 `json:"parent_id"`
	Client               *Client
}

// PermissionOverwrites holds the data of a channel's permission overweries.
type PermissionOverwrites struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Allow int    `json:"allow"`
	Deny  int    `json:"deny"`
	Role  Role
	User  User
}

// 9 Permission Overwrites
