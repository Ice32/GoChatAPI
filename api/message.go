package api

import "bytes"

type Message struct {
	Type EventType
	Data interface{}
}

type EventType string

const (
	ChannelAdd        EventType = "ChannelAdd"
	ChannelsSubscribe EventType = "ChannelSubscribe"
	Error             EventType = "Error"
)

func (et EventType) String() string {
	names := map[EventType]string{
		ChannelAdd:        "ChannelAdd",
		ChannelsSubscribe: "ChannelSubscribe",
		Error:             "Error",
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

func NewChannelsMessage(data interface{}) *Message {
	return newMessage(ChannelAdd, data)
}
func NewErrorMessage(data string) *Message {
	return newMessage(Error, data)
}
