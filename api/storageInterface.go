package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/api/types"
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"bitbucket.org/KenanSelimovic/GoChatServer/storage"
	"github.com/mitchellh/mapstructure"
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

func (si StorageInterface) GetMessages(channelId string, send chan types.ChatMessage, errorChannel chan string, stopChannel chan bool) {
	storageInstance := storage.NewStorage(si.dbConnection)

	messagesChannel := make(chan interface{})
	go func() {
		err := storageInstance.GetOnChangeWithOrderAndFilter(
			"messages",
			messagesChannel,
			"createdAt",
			storage.ASC,
			"channelId",
			channelId,
		)
		if err != nil {
			helpers.LogError(err)
		}
	}()
	var newValue types.ChatMessage

	for {
		select {
		case message := <-messagesChannel:
			if err := mapstructure.Decode(message, &newValue); err != nil {
				helpers.LogError(err)
				errorChannel <- err.Error()
			}
			send <- newValue
		case <-stopChannel:
			return
		}
	}
}
func (si StorageInterface) AddChannel(channel string) error {
	storageInstance := storage.NewStorage(si.dbConnection)
	err := storageInstance.Insert("channels", types.NewChannel{Name: channel})
	return err
}
func (si StorageInterface) AddMessage(message types.NewChatMessage) error {
	storageInstance := storage.NewStorage(si.dbConnection)
	err := storageInstance.Insert("messages", message)
	return err
}
func NewStorageInterface(dbConnection *storage.DbConnection) *StorageInterface {
	return &StorageInterface{
		dbConnection: dbConnection,
	}
}
