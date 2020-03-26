package gocord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sacOO7/gowebsocket"
)

// ClientUser holds the data of the user of the client.
type ClientUser struct {
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
	Locale        string `json:"locale"`
	System        bool   `json:"system"`
	Username      string `json:"username"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Tag           string
}

// Client - Client is the core class, starting point of any bot.
// Embeds {ClientUser}
type Client struct {
	User               ClientUser
	HTTPClient         *http.Client
	Ws                 gowebsocket.Socket
	Users              map[string]User
	Guilds             map[string]Guild
	Channel            map[string]GuildChannel
	URL                map[string]string
	APIVersion         string
	Token              string
	Ready              bool
	Connected          bool
	LastHeartbeatSent  int64
	LastAckHeartbeat   int64
	LastSequenceNumber int
	HeartbeatInterval  chan bool
	SessionID          string
	Debug              func(data string)
}

// Login is where the websocket methods always begin.
func (client *Client) Login(Token string) {
	client.Token = Token
	request, requestError := http.NewRequest("GET", client.URL["base"]+"/v"+client.APIVersion+"/gateway/bot", nil)
	if requestError != nil {
		fmt.Println("Request Error on '/gateway/bot': " + requestError.Error())
	}

	if client.Debug != nil {
		client.Debug("Sending request to /gateway/bot to grab initialization data.")
	}

	request.Header.Add("Authorization", "Bot "+client.Token)
	response, httpError := client.HTTPClient.Do(request)
	if httpError != nil {
		fmt.Println("Http Error on '/gateway/bot': " + httpError.Error())
	}

	defer response.Body.Close()
	body, bodyError := ioutil.ReadAll(response.Body)
	if bodyError != nil {
		fmt.Println("Error ooccured when parsing respnse body of '/gateway/bot': " + bodyError.Error())
	}

	var Response GatewayResponse
	json.Unmarshal([]byte(string(body)), &Response)

	if len(Response.URL) < 5 {
		panic("Invalid token provided.")
	}

	if client.Debug != nil {
		client.Debug("URL:" + Response.URL)
		client.Debug("Shards: " + strconv.Itoa(Response.Shards))
		client.Debug("Remaining Session Limit: " + strconv.Itoa(Response.SessionLimit.Remaining))
		client.Debug("Preparing to connect to the gateway...")
	}

	StartConnection(Response, client)
}

// GatewayResponse holds the /gateway/bot response data.
type GatewayResponse struct {
	URL          string `json:"url"`
	Shards       int    `json:"shards"`
	SessionLimit struct {
		Total      int `json:"total"`
		Remaining  int `json:"remaining"`
		ResetAfter int `json:"reset_after"`
	} `json:"session_start_limit"`
}

// Instantiate sets up properties after they've been received from discord.
func (clientUser *ClientUser) Instantiate() {
	clientUser.Tag = clientUser.Username + "#" + clientUser.Discriminator

	if clientUser.System != true && clientUser.System != false {
		clientUser.System = false
	}
}
