package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"net/http"
)

func StartServer(dbConnection *storage.DbConnection) error {
	router := NewRouter(dbConnection)

	router.Handle(ChannelAdd, addChannel)
	router.Handle(ChannelsSubscribe, subscribeForChannels)

	http.Handle("/", router)

	return http.ListenAndServe(":3183", nil)
}
