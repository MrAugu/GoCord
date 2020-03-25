package gocord

import (
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

}

// WebSocketPayload helps deserialize json websocket events.
type WebSocketPayload struct {
	Op   int `json:"op"`
	Data struct {
	} `json:"d"`
	SequenceNumber int `json:"s"`
	EventName      int `json:"t"`
}
