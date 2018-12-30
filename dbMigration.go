package main

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
)

func migrate(session *r.Session, dbName string, tables []string) {
	migrateDatabase(session, dbName)
	migrateTables(session, tables)
}

func migrateDatabase(session *r.Session, dbName string) {
	cursor, err := r.DBList().Run(session)
	databases := make([]string, 10)
	if err = cursor.All(&databases); err != nil {
		log.Fatal(err)
	}
	if !sliceContainsString(databases, dbName) {
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
		if !sliceContainsString(existingTables, table) {
			_, err = r.TableCreate(table).RunWrite(session)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
