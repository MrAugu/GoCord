package gocord

// Guild - Represents data of a "discord server".
type Guild struct {
	AfkChannelID    string
	AfkTimeout      int
	ApplicationID   string
	Available       bool
	Banner          string
	Channels        map[string]GuildChannel
	Client          *Client
	Deleted         bool
	Description     string
	EmbedChannelID  string
	EmbedEnabled    bool
	Icon            string
	ID              string
	JoinedTimestamp int
	Large           bool
	MemberCount     int
	Members         map[string]string
	MfaLevel        int
	Name            string
	OwnerID         string
	Partenered      bool
	Resgion         string
	SystemChannelID string
	Verified        bool
}
