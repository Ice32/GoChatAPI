package api

import (
	"github.com/mitchellh/mapstructure"
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
