package gocord

import "net/http"

// Client represents the starting point of any bot.
type Client struct {
	Channels       ChannelManager
	ReadyTimestamp int
	Token          string
	Uptime         int
	HTTPClient     *http.Client
}

// Login - loggs in the client.
func (client *Client) Login(token string) int {
	client.Token = token
	return 0
}

// InstantiateClient instantiates a client instance with a token.
func InstantiateClient() Client {
	return Client{HTTPClient: &http.Client{}}
}
