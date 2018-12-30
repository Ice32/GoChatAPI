package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type Client struct {
	channel      chan *Message
	stop         chan bool
	socket       *websocket.Conn
	dbConnection *r.Session
}

func (c *Client) forwardFromChannelToSocket() {
	for mess := range c.channel {
		if err := c.socket.WriteJSON(mess); err != nil {
			c.stop <- true
			break
		}
	}
	if err := c.socket.Close(); err != nil {
		fmt.Println(err)
	}
}

func NewClient(socket *websocket.Conn, dbConnection *r.Session) *Client {
	return &Client{
		socket:       socket,
		channel:      make(chan *Message),
		stop:         make(chan bool),
		dbConnection: dbConnection,
	}
}
