// Package database takes care of preparing the database for use by creating
// tables for all database-stored Polly objects. It also provides convenience DB
// and Tx structs that provide common operations on the database such as CRUD
// operations on database-stored Polly objects.
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// Register the pq postgres database driver.
	_ "github.com/lib/pq"
)

// DB is a superset of sqlx.DB that provides common Polly operations.
type DB struct {
	*sqlx.DB
}

// Tx is a superset of sqlx.Tx that provides common Polly operations.
type Tx struct {
	*sqlx.Tx
}

// Config is a struct used to pass along the database configuration.
type Config struct {
	DBName   string
	User     string
	Password string
	SSLMode  string
}

var schema = `
	CREATE TABLE IF NOT EXISTS users (
		id           bigserial PRIMARY KEY,
		token        text,
		display_name text,
		device_type  integer,
		device_guid  text,
		profile_pic  text
	);

	CREATE TABLE IF NOT EXISTS polls (
		id                 bigserial PRIMARY KEY,
		creator_id         bigint references users(id),
		creation_date      timestamp,
		closing_date       timestamp,
		last_updated       timestamp,
		sequence_number    integer,
		last_event_user    text,
		last_event_user_id bigint references users(id),
		last_event_title   text,
		last_event_type    int
	);
`

// Connect opens and connects to a database.
func Connect(config *Config) (*DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("user=%s password=%q dbname=%s sslmode=%s", config.User,
			config.Password, config.DBName, config.SSLMode))
	return &DB{db}, err
}

// Initialize creates all necessary tables if they do not yet exist.
func (db *DB) Initialize() error {
	_, err := db.Exec(schema)
	return err
}

// Begin starts a transaction and returns a *database.Tx instead of a *sql.Tx.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
