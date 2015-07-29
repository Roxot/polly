package database

import (
	"database/sql"
	"fmt"
	"polly"

	_ "polly/internal/github.com/lib/pq"
	"polly/internal/gopkg.in/gorp.v1"
)

const (
	cSSLMode = "disable"
)

type Database struct {
	mapping gorp.DbMap
}

type DBConfig struct {
	Name     string
	User     string
	Password string
}

func NewDatabase(dbConfig *DBConfig) (*Database, error) {
	db := Database{}

	// open the given postgres database
	sqlDB, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Name,
		cSSLMode))

	// return any errors
	if err != nil {
		return &db, err
	}

	// add the tables used, don't yet create them
	db.mapping = gorp.DbMap{Db: sqlDB, Dialect: gorp.PostgresDialect{}}
	db.mapping.AddTableWithName(polly.PrivateUser{}, cUserTableName).
		SetKeys(false, cPK)
	db.mapping.AddTableWithName(polly.Poll{}, cPollTableName).
		SetKeys(true, cPK)
	db.mapping.AddTableWithName(polly.Question{}, cQuestionTableName).
		SetKeys(true, cPK)
	db.mapping.AddTableWithName(polly.Option{}, cOptionTableName).
		SetKeys(true, cPK)
	db.mapping.AddTableWithName(polly.Vote{}, cVoteTableName).
		SetKeys(true, cPK)
	db.mapping.AddTableWithName(polly.Participant{}, cParticipantTableName).
		SetKeys(true, cPK)

	return &db, nil
}

func (db *Database) CreateTablesIfNotExists() error {
	return db.mapping.CreateTablesIfNotExists()
}

func (db *Database) DropTablesIfExists() error {
	return db.mapping.DropTablesIfExists()
}

func (db *Database) Begin() (*gorp.Transaction, error) {
	return db.mapping.Begin()
}

func (db *Database) Close() {
	db.mapping.Db.Close()
}
