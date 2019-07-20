package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client type for sending messages to the server and sending messages from the server to the client
type Client struct {
	name      string
	webSocket *websocket.Conn
	server    *Server
	// buffered channel to send messages to the client
	send chan *OutputMessage
}

//NewClient create new client
func NewClient(clientName string, webSocket *websocket.Conn, server *Server) *Client {
	if webSocket == nil {
		panic("websocket should be not nil")
	}
	if server == nil {
		panic("server should be not nil")
	}

	return &Client{
		clientName,
		webSocket,
		server,
		make(chan *OutputMessage, 25)}
}

//Listen function to listen for incoming and outgoing messages
func (cl *Client) Listen() {
	go cl.readMessage()
	cl.writeMessage()
}

//readMessage listen, sign and send an incoming message
// to the server channel for broadcast to other clients
func (cl *Client) readMessage() {
	defer func() {
		cl.server.removeClient <- cl
		cl.webSocket.Close()
	}()

	for {
		var msg InputMessage
		err := cl.webSocket.ReadJSON(&msg)
		if err != nil {
			log.Printf("client %s %v", cl.name, err)
			break
		}
		outMsg := OutputMessage{cl.name, msg.Body}
		cl.server.broadcast <- &outMsg
	}
}

// writeMessage send message to client
func (cl *Client) writeMessage() {
	for {
		msg := <-cl.send

		err := cl.webSocket.WriteJSON(msg)
		if err != nil {
			log.Printf("client %s %v", cl.name, err)
			break
		}
	}
}
