package gocord

import "net/http"

// SpawnClient returns a freshly initialized empty client.
func SpawnClient() Client {
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
	return initializedClient
}
