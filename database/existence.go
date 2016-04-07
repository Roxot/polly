package database

// import (
// 	"fmt"

// 	_ "github.com/lib/pq"
// 	"gopkg.in/gorp.v1"
// )

// // TODO check when this error occcurs, maybe we could just return a bool
// func (db *Database) ExistsParticipant(userID, pollID int64) (bool, error) {

// 	count, err := db.mapping.SelectInt(fmt.Sprintf(
// 		"select count(1) from %s where %s=$1 and %s=$2;",
// 		cParticipantTableName, cUserID, cPollID), userID, pollID)
// 	if err != nil {
// 		return false, err
// 	}

// 	return (count == 1), nil
// }

// func ExistsParticipantTX(userID, pollID int64, tx *gorp.Transaction) (bool,
// 	error) {

// 	count, err := tx.SelectInt(fmt.Sprintf(
// 		"select count(1) from %s where %s=$1 and %s=$2;",
// 		cParticipantTableName, cUserID, cPollID), userID, pollID)
// 	if err != nil {
// 		return false, err
// 	}

// 	return (count == 1), nil
// }
