package gocord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sacOO7/gowebsocket"
)

// ClientUser - Embeds {User} and holds client's user data.
type ClientUser struct {
	User
	Verified   bool
	MfaEnabled bool
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
	Debug              bool
	LastHeartbeatSent  int64
	LastAckHeartbeat   int64
	LastSequenceNumber int
}

// Login is where the websocket methods always begin.
func (client *Client) Login(Token string) {
	client.Token = Token
	request, requestError := http.NewRequest("GET", client.URL["base"]+"/v"+client.APIVersion+"/gateway/bot", nil)
	if requestError != nil {
		fmt.Println("Request Error on '/gateway/bot': " + requestError.Error())
	}

	if client.Debug == true {
		fmt.Println("Sending request to /gateway/bot to grab initialization data.")
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

	if client.Debug == true {
		fmt.Println("URL:", Response.URL)
		fmt.Println("Shards:", Response.Shards)
		fmt.Println("Remaining Session Limit:", Response.SessionLimit.Remaining)
		fmt.Println("Preparing to connect to the gateway...")
	}

	StartConnection(Response, client)
}

// GatewayResponse holds the /gateway/bot response data.
type GatewayResponse struct {
	URL          string  `json:"url"`
	Shards       float64 `json:"shards"`
	SessionLimit struct {
		Total      int `json:"total"`
		Remaining  int `json:"remaining"`
		ResetAfter int `json:"reset_after"`
	} `json:"session_start_limit"`
}
