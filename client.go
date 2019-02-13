package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type client struct {
	user *ChatUser

	conversation *Conversation

	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan []byte

	// room is the room this client is chatting in.
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, text, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		msg := c.conversation.addMessage(c.user.Id, string(text))
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return
		}

		c.room.forward <- msgBytes
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {

		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
