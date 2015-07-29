package database

import (
	"fmt"

	_ "polly/internal/github.com/lib/pq"
)

// TODO check when this error occcurs, maybe we could just return a bool
func (db *Database) ExistsParticipant(userID int64, pollID int64) (bool,
	error) {

	count, err := db.mapping.SelectInt(fmt.Sprintf(
		"select count(1) from %s where %s=$1 and %s=$2;",
		cParticipantTableName, cUserID, cPollID), userID, pollID)
	if err != nil {
		return false, err
	}

	return (count == 1), nil
}
