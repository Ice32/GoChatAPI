package types

import "time"

type ChatMessage struct {
	Id        string    `rethinkdb:"id"`
	Text      string    `rethinkdb:"text"`
	ChannelId string    `rethinkdb:"channelId"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
