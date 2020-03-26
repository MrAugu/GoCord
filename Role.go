package gocord

// Role - Holds the data of a role.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
	Client      *Client
}

// Instantiate instantiates a Role structure.
func (role *Role) Instantiate(client *Client) {
	role.Client = client
}
