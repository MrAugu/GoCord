package gocord

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sacOO7/gowebsocket"
)

// StartConnection - initializes a new WebSocket connection.
func StartConnection(Response GatewayResponse, client *Client) {
	client.Ws = gowebsocket.New(Response.URL)
	client.Ws.OnConnected = func(socket gowebsocket.Socket) {
		connected(socket, client)
	}
	client.Ws.OnConnectError = func(err error, socket gowebsocket.Socket) {
		connectionError(err, socket, client)
	}
	client.Ws.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		disconnected(err, socket, client)
	}
	client.Ws.OnTextMessage = func(text string, socket gowebsocket.Socket) {
		event(text, socket, client)
	}
	client.Ws.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		fmt.Println("Unexpectedly recieved binary data ", data)
	}
	if client.Debug != nil {
		client.Debug("Connecting to the gateway...")
	}

	client.Ws.Connect()
}

func connected(socket gowebsocket.Socket, client *Client) {
	client.Connected = true
	if client.Debug != nil {
		client.Debug("Websocket connection established.")
	}
}

func connectionError(err error, socket gowebsocket.Socket, client *Client) {
	if client.Debug != nil {
		client.Debug("Websocket connection error: " + err.Error())
	}
}

func disconnected(err error, socket gowebsocket.Socket, client *Client) {
	// Much more logic needed to this.
	if client.Debug != nil {
		client.Debug("Disconnected from the socket.")
	}
}

func event(message string, socket gowebsocket.Socket, client *Client) {
	var payload WebSocketPayload
	json.Unmarshal([]byte(message), &payload)

	if payload.SequenceNumber != 0 {
		client.LastSequenceNumber = payload.SequenceNumber
	}

	if payload.Op == 10 {
		initializeHeartbeat(payload.Data.HeartbeatInterval, client, socket)
		if client.Ready != true {
			var payload string
			var rawPayload identifyGatewayPacket = identifyGatewayPacket{Op: 2, Event: "IDENTIFY", Data: identifyPacket{Token: client.Token, Properties: identifyProps{Browser: "gocord", Os: "windows", Device: "gocord"}, Compress: false, LargeThreshold: 250, Presence: identifyPresence{Afk: false, Status: "online"}}}
			rawIdentifyPayload, _ := json.Marshal(rawPayload)
			payload = string(rawIdentifyPayload)
			if client.Debug != nil {
				client.Debug("Sending the IDENTIFY payload.")
			}
			socket.SendText(payload)
		}
	}

	if payload.Op == 11 {
		client.LastAckHeartbeat = int64(time.Now().Unix())
		if client.Debug != nil {
			client.Debug("Heartbeat acknowledged.")
		}
	}

	if payload.Event == "READY" {
		var BotUser User = payload.Data.User
		client.User = ClientUser{Avatar: BotUser.Avatar, Bot: BotUser.Bot, Discriminator: BotUser.Discriminator, ID: BotUser.ID, Locale: BotUser.Locale, System: BotUser.System, Username: BotUser.Username, MfaEnabled: BotUser.MfaEnabled}
		client.User.Instantiate()

		client.SessionID = payload.Data.SessionID
		client.InitialGuilds = payload.Data.Guilds
		if client.OnReady != nil {
			client.OnReady()
		}

		if client.Debug != nil {
			client.Debug("Received session data, waiting to receive guilds...")
		}
	}

	if payload.Event == "GUILD_CREATE" && client.Ready != true {
		fmt.Println(payload.Data.Name)
	}
}

func initializeHeartbeat(interval int, client *Client, socket gowebsocket.Socket) {
	if client.Debug != nil {
		client.Debug("Setting heartbeat interval at " + strconv.Itoa(interval) + "ms.")
	}

	var channel chan bool = SetInterval(func() {
		var payload string
		packet := HeartbeatPacket{Op: 1, Seq: client.LastSequenceNumber}
		rawPayload, err := json.Marshal(packet)
		if err != nil {
			fmt.Println(err)
		}
		payload = string(rawPayload)
		if client.Debug != nil {
			client.Debug("Sending a heartbeat.")
		}
		socket.SendText(payload)
		client.LastHeartbeatSent = int64(time.Now().Unix())
	}, interval, true)

	client.HeartbeatInterval = channel

	var intervalChan chan bool
	intervalChan = SetInterval(func() {
		if client.LastHeartbeatSent-client.LastAckHeartbeat > 5000 {
			if client.Debug != nil {
				client.Debug("Last heartbeat acknowledged was over 5 seconds ago. Closing...")
			}
			socket.Close()
			close(client.HeartbeatInterval)
			close(intervalChan)
		}
	}, 5000, true)
}

// WebSocketPayload helps deserialize json websocket events.
type WebSocketPayload struct {
	Op   int `json:"op"`
	Data struct {
		HeartbeatInterval           int                `json:"heartbeat_interval"`
		SessionID                   string             `json:"session_id"`
		GatewayVersion              int                `json:"v"`
		User                        User               `json:"user"`
		PrivateChannels             []string           `json:"private_channels"`
		Guilds                      []UnavailableGuild `json:"guilds"`
		ID                          string             `json:"id"`
		Name                        string             `json:"name"`
		Icon                        string             `json:"icon"`
		Splash                      string             `json:"splash"`
		DiscoverySplash             string             `json:"discovery_splash"`
		OwnerID                     string             `json:"owner_id"`
		Permissions                 int                `json:"permissions"`
		Region                      string             `json:"region"`
		AfkChannelID                string             `json:"afk_channel_id"`
		AfkTimeout                  int                `json:"afk_timeout"`
		EmbedEnabled                bool               `json:"embed_enabled"`
		EmbedChannelID              string             `json:"embed_channel_id"`
		VerificationLevel           int                `json:"verification_level"`
		DefaultMessageNotifications int                `json:"default_message_notifications"`
		ExplicitContentFilter       int                `json:"explicit_content_filter"`
		Roles                       []Role             `json:"roles"`
		Emojis                      []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			User struct {
				Avatar        string `json:"avatar"`
				Bot           bool   `json:"bot"`
				Discriminator string `json:"discriminator"`
				ID            string `json:"id"`
				Locale        string `json:"locale"`
				System        bool   `json:"system"`
				Username      string `json:"username"`
				MfaEnabled    bool   `json:"mfa_enabled"`
			} `json:"user"`
			RequireColons bool `json:"require_colons"`
			Managed       bool `json:"managed"`
			Animated      bool `json:"animated"`
			Available     bool `json:"available"`
		} `json:"emojis"`
		Features           []string     `json:"features"`
		MfaLevel           int          `json:"mfa_level"`
		ApplicationID      string       `json:"application_id"`
		WidgetEnabled      bool         `json:"widget_enabled"`
		WidgetChannelID    string       `json:"widget_channel_id"`
		SystemChannelID    string       `json:"system_channel_id"`
		SystemChannelFlags int          `json:"system_channel_flags"`
		RulesChannelID     string       `json:"rules_channel_id"`
		VoiceStates        []VoiceState `json:"voice_states"`
		Members            []struct {
			User     User            `json:"user"`
			Nickname string          `json:"nick"`
			Roles    map[string]Role `json:"roles"`
			Deaf     bool            `json:"deaf"`
			Muted    bool            `json:"mute"`
		} `json:"members"`
		Large                    bool   `json:"large"`
		Unavailable              bool   `json:"unavailable"`
		MemberCount              int    `json:"member_count"`
		MaxPresences             int    `json:"max_presences"`
		MaxMembers               int    `json:"max_members"`
		VanityCode               string `json:"vanity_url_code"`
		Description              string `json:"description"`
		Banner                   string `json:"banner"`
		PremiumTier              int    `json:"premium_tier"`
		PremiumSubscriptionCount int    `json:"premium_subscription_count"`
		PreferredLocale          string `json:"preferred_locale"`
		PublicUpdatesChannelID   string `json:"public_updates_channel_id"`
	} `json:"d"`
	SequenceNumber int    `json:"s"`
	Event          string `json:"t"`
}

// UnavailableGuild holds data of an unavailable guild.
type UnavailableGuild struct {
	ID          string `json:"id"`
	Unavailable bool   `json:"unavailable"`
}

// HeartbeatPacket helps with serialization of the heartbeat packet.
type HeartbeatPacket struct {
	Op  int `json:"op"`
	Seq int `json:"d"`
}

type identifyGatewayPacket struct {
	Op    int            `json:"op"`
	Event string         `json:"t"`
	Data  identifyPacket `json:"d"`
}

type identifyPacket struct {
	Token          string           `json:"token"`
	Properties     identifyProps    `json:"properties"`
	Compress       bool             `json:"compress"`
	LargeThreshold int              `json:"large_threshold"`
	Presence       identifyPresence `json:"presence"`
}

type identifyProps struct {
	Os      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type identifyPresence struct {
	Status string `json:"status"`
	Afk    bool   `json:"afk"`
}
