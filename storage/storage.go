package storage

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type Storage struct {
	dbSession *r.Session
}

func (s Storage) GetOnChange(tableName string, receiver chan interface{}) error {
	cursor, err := r.Table(tableName).Changes(r.ChangesOpts{IncludeInitial: true}).Run(s.dbSession)

	if err != nil {
		helpers.LogError(err)
		return err
	}
	defer cursor.Close()

	var value r.ChangeResponse
	for cursor.Next(&value) {
		receiver <- value.NewValue
	}
	return nil
}
func (s Storage) Insert(tableName string, data interface{}) error {
	_, err := r.Table(tableName).Insert(data).RunWrite(s.dbSession)
	return err
}

func NewStorage(session *r.Session) *Storage {
	return &Storage{
		dbSession: session,
	}
}

type DbConnection = r.Session

func NewDbConnection(dbName string) (*DbConnection, error) {
	return r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbName,
	})
}
