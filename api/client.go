package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"github.com/gorilla/websocket"
)

type Client struct {
	channel             chan *Message
	stop                chan bool
	messagesStopChannel chan bool
	socket              *websocket.Conn
	dbConnection        *storage.DbConnection
}

func (c *Client) forwardFromChannelToSocket() {
	for mess := range c.channel {
		if err := c.socket.WriteJSON(mess); err != nil {
			c.stop <- true
			break
		}
	}
	if err := c.socket.Close(); err != nil {
		helpers.LogError(err)
	}
}

func NewClient(socket *websocket.Conn, dbConnection *storage.DbConnection) *Client {
	return &Client{
		socket:              socket,
		channel:             make(chan *Message),
		messagesStopChannel: make(chan bool),
		stop:                make(chan bool),
		dbConnection:        dbConnection,
	}
}
