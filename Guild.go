package gocord

// Guild - Represents data of a "discord server".
type Guild struct {
	ID                          string                 `json:"id"`
	Name                        string                 `json:"name"`
	Icon                        string                 `json:"icon"`
	Splash                      string                 `json:"splash"`
	DiscoverySplash             string                 `json:"discovery_splash"`
	OwnerID                     string                 `json:"owner_id"`
	Permissions                 int                    `json:"permissions"`
	Region                      string                 `json:"region"`
	AfkChannelID                string                 `json:"afk_channel_id"`
	AfkTimeout                  int                    `json:"afk_timeout"`
	EmbedEnabled                bool                   `json:"embed_enabled"`
	EmbedChannelID              string                 `json:"embed_channel_id"`
	VerificationLevel           int                    `json:"verification_level"`
	DefaultMessageNotifications int                    `json:"default_message_notifications"`
	ExplicitContentFilter       int                    `json:"explicit_content_filter"`
	Roles                       map[string]Role        `json:"roles"`
	Emojis                      map[string]Emoji       `json:"emojis"`
	Features                    []string               `json:"features"`
	MfaLevel                    int                    `json:"mfa_level"`
	ApplicationID               string                 `json:"application_id"`
	WidgetEnabled               bool                   `json:"widget_enabled"`
	WidgetChannelID             string                 `json:"widget_channel_id"`
	SystemChannelID             string                 `json:"system_channel_id"`
	SystemChannelFlags          int                    `json:"system_channel_flags"`
	RulesChannelID              string                 `json:"rules_channel_id"`
	VoiceStates                 map[string]VoiceState  `json:"voice_states"`
	Members                     map[string]GuildMember `json:"members"`
	Large                       bool                   `json:"large"`
	Unavailable                 bool                   `json:"unavailable"`
	MemberCount                 int                    `json:"member_count"`
	MaxPresences                int                    `json:"max_presences"`
	MaxMembers                  int                    `json:"max_members"`
	VanityCode                  string                 `json:"vanity_url_code"`
	Description                 string                 `json:"description"`
	Banner                      string                 `json:"banner"`
	PremiumTier                 int                    `json:"premium_tier"`
	PremiumSubscriptionCount    int                    `json:"premium_subscription_count"`
	PreferredLocale             string                 `json:"preferred_locale"`
	PublicUpdatesChannelID      string                 `json:"public_updates_channel_id"`
	Client                      *Client
}

// Instantiate instantiates a Guild structure.
func (guild *Guild) Instantiate(client *Client) {
	guild.Client = client
}

// 20 Roles - Done
// Emojis - Done
// Features - Done

// 30 VoiceStates
// Members
// Channels
// Presences
