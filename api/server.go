package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"net/http"
)

func StartServer(dbConnection *storage.DbConnection) error {
	router := NewRouter(dbConnection)

	router.Handle(ChannelAdd, addChannel)
	router.Handle(ChannelsSubscribe, subscribeForChannels)
	router.Handle(MessageAdd, addMessage)
	router.Handle(MessageSubscribe, subscribeForMessages)
	router.Handle(MessageUnsubscribe, unsubscribeForMessages)

	http.Handle("/", router)

	return http.ListenAndServe(":3183", nil)
}
