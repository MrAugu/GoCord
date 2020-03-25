package gocord

// Channel - The base channel which holds the channel data, any other channels are derrived from this struct.
type Channel struct {
	Client  *ClientUser
	deleted bool
	ID      string
	Type    string
}

// GuildChannel - Represents a channel of a guild.
type GuildChannel struct {
	Channel
	deletable   bool
	Guild       *Guild
	ID          string
	Name        string
	ParentID    string
	Position    int
	RawPosition int
	Type        string
	Viewable    string
}
