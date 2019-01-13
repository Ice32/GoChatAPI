package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/api/types"
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"github.com/mitchellh/mapstructure"
	"time"
)

type StorageInterface struct {
	dbConnection *storage.DbConnection
}

func (si StorageInterface) GetChannels(send chan types.Channel, errorChannel chan string) {
	storageInstance := storage.NewStorage(si.dbConnection)

	channelChannel := make(chan interface{})
	go storageInstance.GetOnChange("channels", channelChannel)

	var newValue map[string]string

	for channel := range channelChannel {
		if err := mapstructure.Decode(channel, &newValue); err != nil {
			helpers.LogError(err)
			errorChannel <- err.Error()
		}

		send <- types.Channel{
			Name: newValue["name"],
			Id:   newValue["id"],
		}
	}
}

func (si StorageInterface) GetMessages(send chan string, errorChannel chan string) {
	storageInstance := storage.NewStorage(si.dbConnection)

	messagesChannel := make(chan interface{})
	go storageInstance.GetOnChangeWithOrder(
		"messages",
		messagesChannel,
		"createdAt",
		storage.ASC,
	)
	type newValueStruct struct {
		Text      string
		CreatedAt time.Time
	}
	var newValue newValueStruct

	for message := range messagesChannel {
		if err := mapstructure.Decode(message, &newValue); err != nil {
			helpers.LogError(err)
			errorChannel <- err.Error()
		}
		send <- newValue.Text
	}
}
func (si StorageInterface) AddChannel(channel string) error {
	storageInstance := storage.NewStorage(si.dbConnection)
	err := storageInstance.Insert("channels", types.Channel{Name: channel})
	return err
}
func (si StorageInterface) AddMessage(message ChatMessage) error {
	storageInstance := storage.NewStorage(si.dbConnection)
	err := storageInstance.Insert("messages", message)
	return err
}
func NewStorageInterface(dbConnection *storage.DbConnection) *StorageInterface {
	return &StorageInterface{
		dbConnection: dbConnection,
	}
}
