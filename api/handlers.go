package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/api/types"
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SocketHandler func(client *Client, data interface{})

func subscribeForChannels(client *Client, data interface{}) {
	channelsChannel := make(chan types.Channel)
	errorsChannel := make(chan string)
	go NewStorageInterface(client.dbConnection).GetChannels(channelsChannel, errorsChannel)

	for {
		select {
		case channel := <-channelsChannel:
			message := NewChannelsMessage([]types.Channel{channel})
			client.channel <- message
		case error := <-errorsChannel:
			message := NewErrorMessage(error)
			client.channel <- message
		case <-client.stop:
			return
		}
	}
}

func addChannel(client *Client, data interface{}) {
	var channelName string
	if err := mapstructure.Decode(data, &channelName); err != nil {
		client.channel <- NewErrorMessage(err.Error())
	}
	err := NewStorageInterface(client.dbConnection).AddChannel(channelName)
	if err != nil {
		client.channel <- NewErrorMessage(err.Error())
	}
}

func subscribeForMessages(client *Client, data interface{}) {
	var channelId string
	if err := mapstructure.Decode(data, &channelId); err != nil {
		client.channel <- NewErrorMessage(err.Error())
		return
	}
	messagesChannel := make(chan types.ChatMessage)
	errorsChannel := make(chan string)
	go NewStorageInterface(client.dbConnection).GetMessages(messagesChannel, errorsChannel)

	for {
		select {
		case chatMessage := <-messagesChannel:
			message := NewMessagesMessage([]types.ChatMessage{chatMessage})
			client.channel <- message
		case error := <-errorsChannel:
			message := NewErrorMessage(error)
			client.channel <- message
		case <-client.stop:
			return
		}
	}
}

func unsubscribeForMessages(client *Client, data interface{}) {
	var channelId string
	if err := mapstructure.Decode(data, &channelId); err != nil {
		client.channel <- NewErrorMessage(err.Error())
		return
	}

	// TODO: implement
}

func addMessage(client *Client, data interface{}) {
	var messageText string
	if err := mapstructure.Decode(data, &messageText); err != nil {
		helpers.LogError(err)
		client.channel <- NewErrorMessage(err.Error())
	}
	err := NewStorageInterface(client.dbConnection).AddMessage(types.ChatMessage{
		Text:      messageText,
		CreatedAt: time.Now(),
	})
	if err != nil {
		client.channel <- NewErrorMessage(err.Error())
	}
}
