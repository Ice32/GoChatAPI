package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/api/types"
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SocketHandler func(client *Client, data interface{})
type NewMessageData struct {
	Text      string
	ChannelId string
}

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
	go NewStorageInterface(client.dbConnection).GetMessages(channelId, messagesChannel, errorsChannel, client.messagesStopChannel)

	for {
		select {
		case chatMessage := <-messagesChannel:
			message := NewMessagesMessage([]types.ChatMessage{chatMessage})
			client.channel <- message
		case error := <-errorsChannel:
			message := NewErrorMessage(error)
			client.channel <- message
		case <-client.stop:
		case <-client.messagesStopChannel:
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
	go func() {
		client.messagesStopChannel <- true
	}()
}

func addMessage(client *Client, data interface{}) {
	var message NewMessageData
	if err := mapstructure.Decode(data, &message); err != nil {
		helpers.LogError(err)
		client.channel <- NewErrorMessage(err.Error())
	}
	err := NewStorageInterface(client.dbConnection).AddMessage(types.NewChatMessage{
		Text:      message.Text,
		ChannelId: message.ChannelId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		helpers.LogError(err)
		client.channel <- NewErrorMessage(err.Error())
	}
}
