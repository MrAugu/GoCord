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
	User       ClientUser
	HTTPClient *http.Client
	WebSocket  gowebsocket.Socket
	Users      map[string]User
	Guilds     map[string]Guild
	Channel    map[string]GuildChannel
	URL        map[string]string
	APIVersion string
	Token      string
}

// Login is where the websocket methods always begin.
func (client *Client) Login(Token string) {
	request, requestError := http.NewRequest("GET", client.URL["base"]+"/v"+client.APIVersion+"/gateway/bot", nil)
	if requestError != nil {
		fmt.Println("Request Error on '/gateway/bot': " + requestError.Error())
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
	fmt.Println("Starting a connection...")
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
