package types

import "time"

type NewChatMessage struct {
	Text      string    `rethinkdb:"text"`
	ChannelId string    `rethinkdb:"channelId"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
