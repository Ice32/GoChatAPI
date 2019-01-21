package types

import "time"

type NewChatMessage struct {
	Text      string    `rethinkdb:"text"`
	ChannelId string    `rethinkdb:"channelId"`
	Author    string    `rethinkdb:"author"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
