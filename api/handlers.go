package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SocketHandler func(client *Client, data interface{})

func subscribeForChannels(client *Client, data interface{}) {
	channelsChannel := make(chan string)
	errorsChannel := make(chan string)
	go NewStorageInterface(client.dbConnection).GetChannels(channelsChannel, errorsChannel)

	for {
		select {
		case channelName := <-channelsChannel:
			message := NewChannelsMessage([]string{channelName})
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
	messagesChannel := make(chan string)
	errorsChannel := make(chan string)
	go NewStorageInterface(client.dbConnection).GetMessages(messagesChannel, errorsChannel)

	for {
		select {
		case messageText := <-messagesChannel:
			message := NewMessagesMessage([]string{messageText})
			client.channel <- message
		case error := <-errorsChannel:
			message := NewErrorMessage(error)
			client.channel <- message
		case <-client.stop:
			return
		}
	}
}

func addMessage(client *Client, data interface{}) {
	var messageText string
	if err := mapstructure.Decode(data, &messageText); err != nil {
		helpers.LogError(err)
		client.channel <- NewErrorMessage(err.Error())
	}
	err := NewStorageInterface(client.dbConnection).AddMessage(ChatMessage{
		Text:      messageText,
		CreatedAt: time.Now(),
	})
	if err != nil {
		client.channel <- NewErrorMessage(err.Error())
	}
}
