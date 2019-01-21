package types

import "time"

type ChatMessage struct {
	Id        string    `rethinkdb:"id"`
	Text      string    `rethinkdb:"text"`
	Author    string    `rethinkdb:"author"`
	ChannelId string    `rethinkdb:"channelId"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
