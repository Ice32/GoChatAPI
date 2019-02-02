package storage

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"os"
)

type Storage struct {
	dbSession *r.Session
}
type OrderDirection int

const (
	ASC OrderDirection = iota
	DESC
)

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
func (s Storage) GetOnChangeWithOrder(tableName string, receiver chan interface{}, orderColumn string, orderDirection OrderDirection) error {
	var query r.Term
	if orderDirection == ASC {
		query = r.Table(tableName).OrderBy(r.OrderByOpts{Index: r.Asc(orderColumn)})
	} else {
		query = r.Table(tableName).OrderBy(r.OrderByOpts{Index: r.Desc(orderColumn)})
	}
	cursor, err := query.Limit(1000).Changes(r.ChangesOpts{IncludeInitial: true}).Run(s.dbSession)

	if err != nil {
		helpers.LogError(err)
		return err
	}
	defer func() {
		err = cursor.Close()
		if err != nil {
			helpers.LogError(err)
		}
	}()

	var value r.ChangeResponse
	for cursor.Next(&value) {
		receiver <- value.NewValue
	}
	return nil
}
func (s Storage) GetOnChangeWithOrderAndFilter(tableName string, receiver chan interface{}, orderColumn string, orderDirection OrderDirection, filterField string, filterValue interface{}) error {
	query := r.Table(tableName)
	if orderDirection == ASC {
		query = query.OrderBy(r.OrderByOpts{Index: r.Asc(orderColumn)})
	} else {
		query = query.OrderBy(r.OrderByOpts{Index: r.Desc(orderColumn)})
	}
	cursor, err := query.Filter(r.Row.Field(filterField).Eq(filterValue)).Limit(1000).Changes(r.ChangesOpts{IncludeInitial: true}).Run(s.dbSession)

	if err != nil {
		helpers.LogError(err)
		return err
	}
	defer func() {
		err = cursor.Close()
		if err != nil {
			helpers.LogError(err)
		}
	}()

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
	dbUrl := os.Getenv("DB_HOST")
	if dbUrl == "" {
		dbUrl = "localhost:28015"
	}
	return r.Connect(r.ConnectOpts{
		Address:  dbUrl,
		Database: dbName,
	})
}
