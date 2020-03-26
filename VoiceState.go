package gocord

// VoiceState - Holds a voice state's data.
type VoiceState struct {
	GuildID    string `json:"guild_id"`
	ChannelID  string `json:"channel_id"`
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	Deaf       bool   `json:"deaf"`
	Mute       bool   `json:"mute"`
	SelfDeaf   bool   `json:"self_deaf"`
	SelfMute   bool   `json:"self_mute"`
	SelfStream bool   `json:"self_stream"`
	Suppress   bool   `json:"suppress"`
}

// 8 Member
