package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"github.com/mitchellh/mapstructure"
)

type StorageInterface struct {
	dbConnection *storage.DbConnection
}

func (si StorageInterface) GetChannels(send chan string, errorChannel chan string) {
	storageInstance := storage.NewStorage(si.dbConnection)

	channelChannel := make(chan interface{})
	go storageInstance.GetOnChange("channels", channelChannel)

	var newValue map[string]string

	for channel := range channelChannel {
		if err := mapstructure.Decode(channel, &newValue); err != nil {
			helpers.LogError(err)
			errorChannel <- err.Error()
		}

		send <- newValue["name"]
	}
}
func (si StorageInterface) AddChannel(channel string) error {
	storageInstance := storage.NewStorage(si.dbConnection)
	err := storageInstance.Insert("channels", Channel{Name: channel})
	return err
}
func NewStorageInterface(dbConnection *storage.DbConnection) *StorageInterface {
	return &StorageInterface{
		dbConnection: dbConnection,
	}
}
