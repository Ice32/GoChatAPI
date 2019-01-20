package api

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/api/types"
	"bytes"
)

type Message struct {
	Type EventType
	Data interface{}
}

type EventType string

const (
	ChannelAdd         EventType = "ChannelAdd"
	ChannelsSubscribe  EventType = "ChannelSubscribe"
	MessageAdd         EventType = "MessageAdd"
	MessageSubscribe   EventType = "MessageSubscribe"
	MessageUnsubscribe EventType = "MessageUnsubscribe"
	Error              EventType = "Error"
)

func (et EventType) String() string {
	names := map[EventType]string{
		ChannelAdd:         "ChannelAdd",
		ChannelsSubscribe:  "ChannelSubscribe",
		MessageAdd:         "MessageAdd",
		MessageSubscribe:   "MessageSubscribe",
		MessageUnsubscribe: "MessageUnsubscribe",
		Error:              "Error",
	}

	return names[et]
}

func (et EventType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(et.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func newMessage(eventType EventType, data interface{}) *Message {
	return &Message{
		Type: eventType,
		Data: data,
	}
}

func NewChannelsMessage(data []types.Channel) *Message {
	return newMessage(ChannelAdd, data)
}
func NewErrorMessage(data string) *Message {
	return newMessage(Error, data)
}
func NewMessagesMessage(data interface{}) *Message {
	return newMessage(MessageAdd, data)
}
