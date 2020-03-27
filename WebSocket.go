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

		if len(client.InitialGuilds) < 1 {
			if client.Debug != nil {
				client.Debug("Client has no guilds to receive. Marking it as ready.")
			}
			client.Ready = true
			if client.OnReady != nil {
				client.OnReady()
			}
		}
	}

	if payload.Event == "GUILD_CREATE" && client.Ready != true {
		g := payload.Data
		var createdGuild Guild = Guild{ID: g.ID, Name: g.Name, Icon: g.Icon, Splash: g.Splash, DiscoverySplash: g.DiscoverySplash, OwnerID: g.OwnerID, Permissions: g.Permissions, Region: g.Region, AfkChannelID: g.AfkChannelID, AfkTimeout: g.AfkTimeout, EmbedEnabled: g.EmbedEnabled, EmbedChannelID: g.EmbedChannelID, VerificationLevel: g.VerificationLevel, DefaultMessageNotifications: g.DefaultMessageNotifications, ExplicitContentFilter: g.ExplicitContentFilter, Features: g.Features, MfaLevel: g.MfaLevel, ApplicationID: g.ApplicationID, WidgetEnabled: g.WidgetEnabled, WidgetChannelID: g.WidgetChannelID, SystemChannelID: g.SystemChannelID, SystemChannelFlags: g.SystemChannelFlags, RulesChannelID: g.RulesChannelID, Large: g.Large, Unavailable: g.Unavailable, MemberCount: g.MemberCount, MaxPresences: g.MaxPresences, MaxMembers: g.MaxMembers, VanityCode: g.VanityCode, Description: g.Description, Banner: g.Banner, PremiumTier: g.PremiumTier, PremiumSubscriptionCount: g.PremiumSubscriptionCount, PreferredLocale: g.PreferredLocale, PublicUpdatesChannelID: g.PublicUpdatesChannelID}
		guildRoles := make(map[string]Role)
		for _, key := range payload.Data.Roles {
			properRole := Role{ID: key.ID, Name: key.Name, Color: key.Color, Hoist: key.Hoist, Position: key.Position, Permissions: key.Permissions, Managed: key.Managed, Mentionable: key.Mentionable}
			properRole.Instantiate(client)
			guildRoles[key.ID] = properRole
		}
		createdGuild.Roles = guildRoles

		guildMembers := make(map[string]GuildMember)
		for _, key := range payload.Data.Members {
			user := User{Avatar: key.User.Avatar, Bot: key.User.Bot, Discriminator: key.User.Discriminator, ID: key.User.ID, Locale: key.User.Locale, System: key.User.System, Username: key.User.Username, MfaEnabled: key.User.MfaEnabled}
			user.Instantiate(client)
			client.Users[user.ID] = user

			properMember := GuildMember{User: client.Users[user.ID], Nickname: key.Nickname, Deaf: key.Deaf, Muted: key.Muted}
			properMember.Roles = make(map[string]Role)
			for _, roleKey := range key.Roles {
				properMember.Roles[roleKey] = createdGuild.Roles[roleKey]
			}
			guildMembers[user.ID] = properMember
		}
		createdGuild.Members = guildMembers

		guildEmojis := make(map[string]Emoji)
		for _, key := range payload.Data.Emojis {
			properEmoji := Emoji{ID: key.ID, Name: key.Name, User: client.Users[key.User.ID], RequireColons: key.RequireColons}
			properEmoji.Instantiate(client)
			guildEmojis[key.ID] = properEmoji
		}
		createdGuild.Emojis = guildEmojis

		guildChannels := make(map[string]Channel)
		for _, key := range payload.Data.Channels {
			properChannel := Channel{ID: key.ID, GuildID: key.GuildID, Position: key.Position, Name: key.Name, Topic: key.Topic, Nsfw: key.Nsfw, LastMessageID: key.LastMessageID, Bitrate: key.Bitrate, UserLimit: key.UserLimit, Cooldown: key.Cooldown, Icon: key.Icon, OwnerID: key.OwnerID, ApplicationID: key.ApplicationID, ParentID: key.ParentID}

			channelPerms := make([]PermissionOverwrite, 0, len(key.PermissionOverwrites))
			for _, overwrite := range key.PermissionOverwrites {
				properOverwrite := PermissionOverwrite{ID: overwrite.ID, Type: overwrite.Type}
				properOverwrite.Allow = CalculateBitfield(overwrite.Allow)
				properOverwrite.Deny = CalculateBitfield(overwrite.Deny)

				if properOverwrite.Type == "role" {
					properOverwrite.Role = createdGuild.Roles[overwrite.ID]
				} else {
					properOverwrite.User = client.Users[overwrite.ID]
				}

				channelPerms[len(channelPerms)] = properOverwrite
			}

			properChannel.PermissionOverwrites = channelPerms
			guildChannels[properChannel.ID] = properChannel
			client.Channels[properChannel.ID] = properChannel
		}

		createdGuild.Channels = guildChannels

		guildVoiceStates := make(map[string]VoiceState)
		for _, key := range payload.Data.VoiceStates {
			properVoiceState := VoiceState{GuildID: key.GuildID, ChannelID: key.ChannelID, UserID: key.UserID, SessionID: key.SessionID, Deaf: key.Deaf, Mute: key.Mute, SelfDeaf: key.SelfDeaf, SelfMute: key.SelfMute, SelfStream: key.SelfStream, Suppress: key.Suppress}
			properVoiceState.User = client.Users[key.UserID]

			properVoiceState.Instantiate(client)
			guildVoiceStates[key.UserID] = properVoiceState
		}
		createdGuild.VoiceStates = guildVoiceStates
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
		Channels []struct {
			ID                   string `json:"id"`
			Type                 int    `json:"type"`
			GuildID              string `json:"guild_id"`
			Position             int    `json:"position"`
			Name                 string `json:"name"`
			Topic                string `json:"topic"`
			Nsfw                 bool   `json:"nsfw"`
			LastMessageID        string `json:"last_message_id"`
			Bitrate              int    `json:"bitrate"`
			UserLimit            int    `json:"user_limit"`
			Cooldown             int    `json:"rate_limit_per_user"`
			PermissionOverwrites []struct {
				ID    string `json:"id"`
				Type  string `json:"type"`
				Allow int    `json:"allow"`
				Deny  int    `json:"deny"`
			} `json:"permission_overwrites"`
			Recipients []struct {
				Avatar        string `json:"avatar"`
				Bot           bool   `json:"bot"`
				Discriminator string `json:"discriminator"`
				ID            string `json:"id"`
				Locale        string `json:"locale"`
				System        bool   `json:"system"`
				Username      string `json:"username"`
				MfaEnabled    bool   `json:"mfa_enabled"`
				Tag           string
			} `json:"recipients"`
			Icon          string `json:"icon"`
			OwnerID       string `json:"owner_id"`
			ApplicationID string `json:"application_id"`
			ParentID      string `json:"parent_id"`
		}
		Features           []string `json:"features"`
		MfaLevel           int      `json:"mfa_level"`
		ApplicationID      string   `json:"application_id"`
		WidgetEnabled      bool     `json:"widget_enabled"`
		WidgetChannelID    string   `json:"widget_channel_id"`
		SystemChannelID    string   `json:"system_channel_id"`
		SystemChannelFlags int      `json:"system_channel_flags"`
		RulesChannelID     string   `json:"rules_channel_id"`
		VoiceStates        []struct {
			GuildID    string `json:"guild_id"`
			ChannelID  string `json:"channel_id"`
			UserID     string `json:"user_id"`
			SessionID  string `json:"session_id"`
			Deaf       bool   `json:"deaf"`
			Mute       bool   `json:"mute"`
			SelfDeaf   bool   `json:"self_deaf"`
			SelfMute   bool   `json:"self_mute"`
			SelfStream bool   `json:"self_stream"`
			Suppress   bool   `json:"suppress"`
		} `json:"voice_states"`
		Members []struct {
			User struct {
				Avatar        string `json:"avatar"`
				Bot           bool   `json:"bot"`
				Discriminator string `json:"discriminator"`
				ID            string `json:"id"`
				Locale        string `json:"locale"`
				System        bool   `json:"system"`
				Username      string `json:"username"`
				MfaEnabled    bool   `json:"mfa_enabled"`
				Tag           string
			} `json:"user"`
			Nickname string   `json:"nick"`
			Roles    []string `json:"roles"`
			Deaf     bool     `json:"deaf"`
			Muted    bool     `json:"mute"`
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
