package storage

import (
	"bitbucket.org/KenanSelimovic/GoChatServer/helpers"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
)

func Migrate(session *r.Session, dbName string, tables []string, indexes map[string][]string) {
	migrateDatabase(session, dbName)
	migrateTables(session, tables)
	migrateIndexes(session, indexes)
}

func migrateDatabase(session *r.Session, dbName string) {
	cursor, err := r.DBList().Run(session)
	databases := make([]string, 10)
	if err = cursor.All(&databases); err != nil {
		log.Fatal(err)
	}
	if !helpers.SliceContainsString(databases, dbName) {
		_, err = r.DBCreate(dbName).RunWrite(session)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func migrateTables(session *r.Session, tables []string) {
	existingTables := make([]string, 10)
	tablesCursor, err := r.TableList().Run(session)
	if err = tablesCursor.All(&existingTables); err != nil {
		log.Fatal(err)
	}
	for _, table := range tables {
		if !helpers.SliceContainsString(existingTables, table) {
			_, err = r.TableCreate(table).RunWrite(session)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
func migrateIndexes(session *r.Session, indexes map[string][]string) {
	for table, indexes := range indexes {
		for _, index := range indexes {
			_ = r.Table(table).IndexCreate(index).Exec(session)
		}
	}
}
