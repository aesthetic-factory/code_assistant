package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *Database
	once     sync.Once
)

type Database struct {
	db *sql.DB
}

func GetDatabase() *Database {
	return database
}

func NewDatabase(dbPath string) (*Database, error) {
	var err error
	once.Do(func() {
		var db *sql.DB
		db, err = sql.Open("sqlite3", dbPath)
		if err == nil {
			database = &Database{db: db}
		}
	})
	return database, err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateTable(name string, schema string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, schema)
	_, err := d.db.Exec(query)
	return err
}

func (d *Database) DropTable(name string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
	_, err := d.db.Exec(query)
	return err
}

func (d *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
