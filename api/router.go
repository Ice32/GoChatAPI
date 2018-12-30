package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Router struct {
	handlers     map[EventType]SocketHandler
	dbConnection *storage.DbConnection
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	client := NewClient(socket, router.dbConnection)

	go client.forwardFromChannelToSocket()
	router.handleClient(client)
}

func (router *Router) handleClient(client *Client) {
	var message Message

	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler := router.findHandler(message.Type); handler != nil {
			(*handler)(client, message.Data)
		}
	}

	subscribeForChannels(client, nil)
}

func (router *Router) Handle(eventType EventType, socketHandler SocketHandler) {
	router.handlers[eventType] = socketHandler
}

func (router *Router) findHandler(eventType EventType) *SocketHandler {
	if val, found := router.handlers[eventType]; found {
		return &val
	}
	return nil
}

func NewRouter(session *storage.DbConnection) *Router {
	return &Router{
		handlers:     make(map[EventType]SocketHandler),
		dbConnection: session,
	}
}
