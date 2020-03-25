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
}

func connected(socket gowebsocket.Socket, client *Client) {
	fmt.Println("Websocket connection established.")
}
