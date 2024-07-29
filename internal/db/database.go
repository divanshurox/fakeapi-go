package db

import "context"

type DB interface {
	Connect(*Config) (context.CancelFunc, error)
	Get(*Query, interface{}) error
	Insert(*Query) error
	Update(*Query) error
	Close()
}

type Database struct {
	db DB
}

func GetDatabase(db DB) *Database {
	return &Database{db: db}
}

func (d *Database) Connect(config *Config) (context.CancelFunc, error) {
	return d.db.Connect(config)
}

func (d *Database) Insert(query *Query) error {
	return d.db.Insert(query)
}

func (d *Database) Update(query *Query) error {
	return d.db.Update(query)
}

func (d *Database) Get(query *Query, target interface{}) error {
	return d.db.Get(query, target)
}

func (d *Database) Close() {
	d.db.Close()
}
