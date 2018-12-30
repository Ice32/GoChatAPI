package main

import (
	"github.com/mitchellh/mapstructure"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type Storage struct {
	dbSession *r.Session
}

func (s Storage) GetChannels(send chan string, errorChannel chan string) {
	cursor, err := r.Table("channels").Changes(r.ChangesOpts{IncludeInitial: true}).Run(s.dbSession)

	if err != nil {
		logError(err.Error())
		errorChannel <- err.Error()
	}
	defer cursor.Close()

	var value r.ChangeResponse
	for cursor.Next(&value) {
		var newValue map[string]string

		if err = mapstructure.Decode(value.NewValue, &newValue); err != nil {
			logError(err.Error())
			errorChannel <- err.Error()
		}
		send <- newValue["name"]
	}
}
func (s Storage) AddChannel(channel string) error {
	_, err := r.Table("channels").Insert(Channel{Name: channel}).RunWrite(s.dbSession)
	if err != nil {
		logError(err.Error())
	}
	return err
}

func NewStorage(session *r.Session) *Storage {
	return &Storage{
		dbSession: session,
	}
}
