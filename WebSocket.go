package gocord

import (
	"encoding/json"
	"fmt"
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
	if client.Debug == true {
		fmt.Println("Connecting to the gateway...")
	}

	client.Ws.Connect()
}

func connected(socket gowebsocket.Socket, client *Client) {
	client.Connected = true
	if client.Debug == true {
		fmt.Println("Websocket connection established.")
	}
}

func connectionError(err error, socket gowebsocket.Socket, client *Client) {
	if client.Debug == true {
		fmt.Println("Websocket connection error: " + err.Error())
	}
}

func disconnected(err error, socket gowebsocket.Socket, client *Client) {
	// Much more logic needed to this.
	if client.Debug == true {
		fmt.Println("Disconnected from the socket.")
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
			if client.Debug == true {
				fmt.Println("Sending the IDENTIFY payload.")
			}
			socket.SendText(payload)
		}
	}

	if payload.Op == 11 {
		client.LastAckHeartbeat = int64(time.Now().Unix())
	}

	fmt.Println(payload.Event)
}

func initializeHeartbeat(interval int, client *Client, socket gowebsocket.Socket) {
	if client.Debug == true {
		fmt.Println("Setting heartbeat interval at", interval, "ms.")
	}

	var channel chan bool = SetInterval(func() {
		var payload string
		packet := HeartbeatPacket{Op: 1, Seq: client.LastSequenceNumber}
		rawPayload, err := json.Marshal(packet)
		if err != nil {
			fmt.Println(err)
		}
		payload = string(rawPayload)
		if client.Debug == true {
			fmt.Println("Sending a heartbeat.")
		}
		socket.SendText(payload)
		client.LastHeartbeatSent = int64(time.Now().Unix())
	}, interval, true)

	client.HeartbeatInterval = channel

	var intervalChan chan bool
	intervalChan = SetInterval(func() {
		if client.LastHeartbeatSent-client.LastAckHeartbeat > 5000 {
			if client.Debug == true {
				fmt.Println("Last heartbeat acknowledged was over 5 seconds ago. Closing...")
			}
			socket.Close()
			close(client.HeartbeatInterval)
			close(intervalChan)
		}
	}, interval+10000, true)
}

// WebSocketPayload helps deserialize json websocket events.
type WebSocketPayload struct {
	Op   int `json:"op"`
	Data struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
	SequenceNumber int    `json:"s"`
	Event          string `json:"t"`
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
