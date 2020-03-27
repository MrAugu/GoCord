package gocord

import (
	"net/http"
)

// SpawnClient returns a freshly initialized empty client.
func SpawnClient(debug func(text string)) Client {
	var initializedClient Client = Client{HTTPClient: &http.Client{}}
	var urlMap map[string]string = map[string]string{
		"base":   "https://discordapp.com/api",
		"cdn":    "https://cdn.discordapp.com",
		"invite": "https://discord.gg",
	}
	initializedClient.APIVersion = "6"
	initializedClient.URL = urlMap
	initializedClient.Ready = false
	initializedClient.Connected = false
	initializedClient.Debug = debug

	userMap := make(map[string]User)
	guildMap := make(map[string]Guild)
	channelMap := make(map[string]Channel)

	initializedClient.Users = userMap
	initializedClient.Guilds = guildMap
	initializedClient.Channels = channelMap

	return initializedClient
}
