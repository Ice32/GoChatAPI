package main

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
	"net/http"
)

const dbName string = "goChat"

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbName,
	})
	if err != nil {
		log.Fatalln(err)
	}

	migrate(
		session,
		dbName,
		[]string{"channels"},
	)

	router := NewRouter(session)

	router.Handle(ChannelAdd, addChannel)
	router.Handle(ChannelsSubscribe, subscribeForChannels)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3183", nil))
}
