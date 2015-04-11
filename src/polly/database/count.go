package database

import (
	"fmt"

	_ "github.com/lib/pq"
)

func (db *Database) CountPollsForUser(usrId int) int64 {

	count, err := db.dbMap.SelectInt(fmt.Sprintf(
		"select count(*) from %s where %s=$1;", cParticipantTableName, cUserId),
		usrId)
	if err != nil {
		return 0 // TODO when is error nil?
	}

	return count
}
