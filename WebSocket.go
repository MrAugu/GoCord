package gocord

import (
	"encoding/json"
	"fmt"

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
	fmt.Println(payload.SequenceNumber)
	if payload.SequenceNumber != 0 {
		client.LastSequenceNumber = payload.SequenceNumber
	}

	if payload.Op == 10 {
		initializeHeartbeat(payload.Data.HeartbeatInterval, client, socket)
	}
}

// WebSocketPayload helps deserialize json websocket events.
type WebSocketPayload struct {
	Op   int `json:"op"`
	Data struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
	SequenceNumber int `json:"s"`
	EventName      int `json:"t"`
}

func initializeHeartbeat(interval int, client *Client, socket gowebsocket.Socket) {
	if client.Debug == true {
		fmt.Print("Setting heartbeat interval at ", interval)
		fmt.Print("ms.")
	}

	SetInterval(func() {
		var payload string
		packet := HeartbeatPacket{Op: 1, Seq: client.LastSequenceNumber}
		rawPayload, err := json.Marshal(packet)
		if err != nil {
			fmt.Println(err)
		}
		payload = string(rawPayload)
		fmt.Println(payload)
		if client.Debug == true {
			fmt.Print("Sending a heartbeat.")
		}
	}, interval, true)
}

// HeartbeatPacket helps with serialization of the heartbeat packet.
type HeartbeatPacket struct {
	Op  int `json:"op"`
	Seq int `json:"d"`
}
