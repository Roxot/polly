package database

import (
	"database/sql"
	"fmt"
	"polly"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

const (
	cSSLMode = "disable"
)

type Database struct {
	dbMap gorp.DbMap
}

type DbConfig struct {
	DbName       string
	PsqlUser     string
	PsqlUserPass string
}

func New(dbConfig DbConfig) (*Database, error) {
	db := Database{}

	// open the given postgres database
	sqlDb, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.PsqlUser, dbConfig.PsqlUserPass, dbConfig.DbName,
			cSSLMode))

	// return any errors
	if err != nil {
		return &db, err
	}

	// add the tables used, don't yet create them
	db.dbMap = gorp.DbMap{Db: sqlDb, Dialect: gorp.PostgresDialect{}}
	db.dbMap.AddTableWithName(polly.PrivateUser{}, cUserTableName).
		SetKeys(true, cPk).ColMap(cPhoneNumber).SetUnique(true)
	db.dbMap.AddTableWithName(polly.VerToken{},
		cVerificationTokensTableName).SetKeys(true, cPk)
	db.dbMap.AddTableWithName(polly.Poll{}, cPollTableName).
		SetKeys(true, cPk)
	db.dbMap.AddTableWithName(polly.Question{}, cQuestionTableName).
		SetKeys(true, cPk)
	db.dbMap.AddTableWithName(polly.Option{}, cOptionTableName).
		SetKeys(true, cPk)
	db.dbMap.AddTableWithName(polly.Vote{}, cVoteTableName).
		SetKeys(true, cPk)
	db.dbMap.AddTableWithName(polly.Participant{}, cParticipantTableName).
		SetKeys(true, cPk)

	return &db, nil
}

func (db *Database) CreateTablesIfNotExists() error {
	return db.dbMap.CreateTablesIfNotExists()
}

func (db *Database) DropTablesIfExists() error {
	return db.dbMap.DropTablesIfExists()
}

func (db *Database) Close() {
	db.dbMap.Db.Close()
}
