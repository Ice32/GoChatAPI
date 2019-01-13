package api

import "time"

type ChatMessage struct {
	Text      string    `rethinkdb:"text"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
