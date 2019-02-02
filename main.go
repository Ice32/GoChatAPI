package main

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/api"
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"log"
)

const dbName string = "goChat"

func main() {
	dbConnection, err := storage.NewDbConnection(dbName)
	if err != nil {
		log.Fatalln(err)
	}

	storage.Migrate(
		dbConnection,
		dbName,
		[]string{"channels", "messages"},
		map[string][]string{
			"messages": {"createdAt"},
		},
	)

	log.Fatal(api.StartServer(dbConnection))
}
