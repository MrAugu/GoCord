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

	if payload.Op == 10 {
		initializeHeartbeat(payload.Data.HeartbeatInterval, client, socket)
	}

	if payload.Op == 11 {
		client.LastAckHeartbeat = int64(time.Now().Unix())
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

// HeartbeatPacket helps with serialization of the heartbeat packet.
type HeartbeatPacket struct {
	Op  int `json:"op"`
	Seq int `json:"d"`
}
