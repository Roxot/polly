package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)
import "fmt"

const (
	cSSLMode = "disable"
)

type Database struct {
	dbMap *gorp.DbMap
}

type DbConfig struct {
	DbName       string
	PsqlUser     string
	PsqlUserPass string
}

func New(dbConfig DbConfig) (Database, error) {

	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.PsqlUser, dbConfig.PsqlUserPass, dbConfig.DbName,
			cSSLMode))

	if err != nil {
		return Database{}, err
	}

	pollyDb := Database{}
	pollyDb.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	pollyDb.dbMap.AddTableWithName(User{}, cUserTableName).SetKeys(true, cPk).
		ColMap(cPhoneNumber).SetUnique(true)
	pollyDb.dbMap.AddTableWithName(VerificationToken{},
		cVerificationTokensTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Poll{}, cPollTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Question{}, cQuestionTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Option{}, cOptionTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Vote{}, cVoteTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Participant{}, cParticipantTableName).SetKeys(true, cPk)

	return pollyDb, nil
}

func (pollyDb Database) CreateTablesIfNotExists() error {
	return pollyDb.dbMap.CreateTablesIfNotExists()
}

func (pollyDb Database) DropTablesIfExists() error {
	return pollyDb.dbMap.DropTablesIfExists()
}

func (pollyDb Database) Close() {
	pollyDb.dbMap.Db.Close()
}
