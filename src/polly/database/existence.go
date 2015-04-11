package database

import (
	"fmt"

	_ "github.com/lib/pq"
)

// TODO check when this error occcurs, maybe we could just return a bool
func (db *Database) ExistsParticipant(usrId, pollId int) (bool, error) {
	count, err := db.dbMap.SelectInt(fmt.Sprintf(
		"select count(1) from %s where %s=$1 and %s=$2;",
		cParticipantTableName, cUserId, cPollId), usrId, pollId)
	if err != nil {
		return false, err
	}

	return (count == 1), nil
}
