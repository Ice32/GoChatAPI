package types

import "time"

type ChatMessage struct {
	Id        string    `rethinkdb:"id"`
	Text      string    `rethinkdb:"text"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
