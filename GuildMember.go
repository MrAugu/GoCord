package gocord

// GuildMember holds the data of a guild member.
type GuildMember struct {
	User     User            `json:"user"`
	Nickname string          `json:"nick"`
	Roles    map[string]Role `json:"roles"`
	Deaf     bool            `json:"deaf"`
	Muted    bool            `json:"mute"`
}
