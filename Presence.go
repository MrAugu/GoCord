package gocord

// RawPresence used to serialize websocket JSON presence payload.
type RawPresence struct {
	User struct {
		Avatar        string `json:"avatar"`
		Bot           bool   `json:"bot"`
		Discriminator string `json:"discriminator"`
		ID            string `json:"id"`
		Locale        string `json:"locale"`
		System        bool   `json:"system"`
		Username      string `json:"username"`
		MfaEnabled    bool   `json:"mfa_enabled"`
	} `json:"user"`
	Roles   []string     `json:"roles"`
	Game    PresenceGame `json:"game"`
	GuildID string       `json:"guild_id"`
	Status  string       `json:"status"`
	Nick    string       `json:"nick"`
}

// PresenceGame helds data of the game an user is playing.
type PresenceGame struct {
	Name       string `json:"name"`
	Type       int    `json:"type"`
	URL        int    `json:"url"`
	CreatedAt  int    `json:"created_at"`
	Timestamps struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"timestamps"`
	ApplicationID string `json:"application_id"`
	Details       string `json:"details"`
	Status        string `json:"state"`
	Emoji         struct {
		Name     string `json:"name"`
		ID       string `json:"id"`
		Animated bool   `json:"animated"`
	} `json:"emoji"`
	Party struct {
		ID   string `json:"id"`
		Size []int  `json:"size"`
	} `json:"party"`
}

// Presence holds the presence data attached to User.
type Presence struct {
	Status string
	Game   PresenceGame
	User   User
}
