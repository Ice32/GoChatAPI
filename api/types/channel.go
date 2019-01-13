package types

type Channel struct {
	Id   string `rethinkdb:"id"`
	Name string `rethinkdb:"name"`
}
