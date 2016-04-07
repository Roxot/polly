// The database takes care of preparing the database for use by creating tables
// for all database-stored Polly objects. It also provides convenience DB and Tx
// implementations that provide common operations on the database such as CRUD
// operations on database-stored Polly objects.
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// A superset of sql.DB that provides common Polly operations.
type DB struct {
	*sqlx.DB
}

// A superset of sqlx.Tx that provides common Polly operations.
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
