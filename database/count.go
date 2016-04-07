package database

import _ "github.com/lib/pq"

// func (db *Database) CountPollsForUser(userID int64) int64 {

// 	count, err := db.mapping.SelectInt(fmt.Sprintf(
// 		"select count(*) from %s where %s=$1;", cParticipantTableName, cUserID),
// 		userID)
// 	if err != nil {
// 		return 0 // TODO when is error nil?
// 	}

// 	return count
// }
