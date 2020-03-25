package gocord

// User - Holds an user's data.
type User struct {
	Client        *Client
	Avatar        string
	Bot           bool
	Discriminator string
	ID            string
	Locale        string
	System        bool
	Tag           string
	Username      string
}
