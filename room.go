package main

import (
	"github.com/gorilla/websocket"
)

type room struct {
	created bool

	id      int
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom(convId int) *room {
	dRoom := &room{
		id:      convId,
		created: true,
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
	rooms[convId] = dRoom
	return dRoom
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			if len(r.clients) < 1 {
				// Lebski - delete room when its empty
				delete(rooms, r.id)
			}
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func roomExists(conversationId int) (*room, bool) {
	dRoom, ok := rooms[conversationId]
	if !ok {
		return &room{}, false
	}
	return dRoom, true

}
