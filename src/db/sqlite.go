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

// GetDatabase returns the database object.
//
// No parameters.
// Returns a pointer to the Database object.
func GetDatabase() *Database {
	return database
}

// NewDatabase creates a new Database instance.
//
// It takes a dbPath string as a parameter and returns a pointer to Database and an error.
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

// Close closes the database connection.
//
// It returns an error if there was a problem closing the database.
func (d *Database) Close() error {
	return d.db.Close()
}

// CreateTable creates a table in the database.
//
// Parameters:
// - name: the name of the table.
// - schema: the schema definition of the table.
// Returns an error.
func (d *Database) CreateTable(name string, schema string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, schema)
	_, err := d.db.Exec(query)
	return err
}

// DropTable drops a table from the database.
//
// Parameters:
// - name: the name of the table to be dropped.
//
// Returns:
// - error: an error if the table drop operation fails.
func (d *Database) DropTable(name string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
	_, err := d.db.Exec(query)
	return err
}

// Execute executes a SQL query and returns the result.
//
// Parameters:
// - query: the SQL query to execute.
// - args: variadic list of arguments to be passed to the query.
// Returns:
// - sql.Result: the result of the SQL query.
// - error: an error if the query execution fails.
func (d *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Query executes the given SQL query with optional arguments on the Database.
//
// query - the SQL query to be executed
// args - optional arguments for the query
// Returns *sql.Rows and error.
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
