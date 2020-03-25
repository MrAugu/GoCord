package gocord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sacOO7/gowebsocket"
)

// Client represents the starting point of any bot.
type Client struct {
	Channels       ChannelManager
	ReadyTimestamp int
	Token          string
	Uptime         int
	HTTPClient     *http.Client
	BaseURL        string
	CdnURL         string
	InviteURL      string
	APIVersion     string
	WebSocket      gowebsocket.Socket
}

// Login - loggs in the client.
func (client *Client) Login(token string) int {
	client.Token = token
	request, requestError := http.NewRequest("GET", client.BaseURL+"/v"+client.APIVersion+"/gateway/bot", nil)

	if requestError != nil {
		fmt.Println("Request Error on '/gateway/bot': " + requestError.Error())
		return 1
	}

	request.Header.Add("Authorization", "Bot "+client.Token)
	response, httpError := client.HTTPClient.Do(request)

	if httpError != nil {
		fmt.Println("Http Error on '/gateway/bot': " + httpError.Error())
		return 1
	}

	defer response.Body.Close()
	body, bodyError := ioutil.ReadAll(response.Body)

	if bodyError != nil {
		fmt.Println("Error ooccured when parsing respnse body of '/gateway/bot': " + bodyError.Error())
		return 1
	}

	var InitData GatewayInitData
	json.Unmarshal([]byte(string(body)), &InitData)

	client.WebSocket = gowebsocket.New(InitData.URL)
	client.WebSocket.OnConnected = func(socket gowebsocket.Socket) {
		onReadySocketEvent(socket, client)
	}

	client.WebSocket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		onTextMessageSocketEvent(message, socket, client)
	}

	client.WebSocket.Connect()
	return 0
}

// InstantiateClient instantiates a client instance with a token.
func InstantiateClient() Client {
	return Client{HTTPClient: &http.Client{}, BaseURL: "https://discordapp.com/api", CdnURL: "https://cdn.discordapp.com", InviteURL: "InviteURL", APIVersion: "7"}
}

// GatewayInitData - Holds gateway data fetched from /gatway/bot.
type GatewayInitData struct {
	URL          string                   `json:"url"`
	Shards       float64                  `json:"shards"`
	SessionLimit GatewaySessionStartLimit `json:"session_start_limit"`
}

// GatewaySessionStartLimit is a structure that hold the nested object of the initial /gateway/bot response.
type GatewaySessionStartLimit struct {
	Total      int `json:"total"`
	Remaining  int `json:"remaining"`
	ResetAfter int `json:"reset_after"`
}

func onReadySocketEvent(socket gowebsocket.Socket, client *Client) {
	fmt.Println("WebSocket connection established.")
}

func onTextMessageSocketEvent(text string, socket gowebsocket.Socket, client *Client) {
	fmt.Println("Received: " + text)
}
