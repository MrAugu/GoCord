package gocord

// Client represents the starting point of any bot.
type Client struct {
	Channels       ChannelManager
	ReadyTimestamp int
	Token          string
	Uptime         int
}

// LoginClient initializes a connection to discord.
func LoginClient(client Client) {

}
