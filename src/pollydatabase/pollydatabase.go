package pollydatabase

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)
import "fmt"

const (
	cDatabaseName = "pollydb"
	cUserName     = "polly"
	cPassword     = "w01V3s"
	cSSLMode      = "disable"
)

type PollyDatabase struct {
	dbMap *gorp.DbMap
}

func New() (PollyDatabase, error) {

	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", cUserName,
			cPassword, cDatabaseName, cSSLMode))

	if err != nil {
		return PollyDatabase{}, err
	}

	pollyDb := PollyDatabase{}
	pollyDb.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	pollyDb.dbMap.AddTableWithName(User{}, cUserTableName).SetKeys(true, cPk).
		ColMap(cPhoneNumber).SetUnique(true)
	pollyDb.dbMap.AddTableWithName(Poll{}, cPollTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Question{}, cQuestionTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Option{}, cOptionTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Vote{}, cVoteTableName).SetKeys(true, cPk)
	pollyDb.dbMap.AddTableWithName(Participant{}, cParticipantTableName).SetKeys(true, cPk)

	return pollyDb, nil
}

func (pollyDb PollyDatabase) CreateTables() error {
	return pollyDb.dbMap.CreateTablesIfNotExists()
}

func (pollyDb PollyDatabase) DropTables() error {
	return pollyDb.dbMap.DropTablesIfExists()
}

func (pollyDb PollyDatabase) Close() {
	pollyDb.dbMap.Db.Close()
}
