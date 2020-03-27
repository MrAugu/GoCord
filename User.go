package gocord

// User - Holds an user's data.
type User struct {
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
	Locale        string `json:"locale"`
	System        bool   `json:"system"`
	Username      string `json:"username"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Tag           string
	Client        *Client
	Presence      Presence
}

// Instantiate instantiates an User structure.
func (user *User) Instantiate(client *Client) {
	user.Tag = user.Username + "#" + user.Discriminator
	user.Client = client

	if user.System != true && user.System != false {
		user.System = false
	}
}
