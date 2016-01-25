package database

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

func DeleteVotesForUserTX(userID, pollID int64, tx *gorp.Transaction) error {
	_, err := tx.Exec(fmt.Sprintf("delete from %s where %s=$1 and %s=$2;",
		cVoteTableName, cUserID, cPollID), userID, pollID)
	return err
}

func (db *Database) DeleteVoteByIDForUser(voteID, userID int64) error {
	_, err := db.mapping.Exec(fmt.Sprintf(
		"delete from %s where %s=$1 and %s=$2;", cVoteTableName, cID,
		cUserID), voteID, userID)
	return err
}

func (db *Database) DeleteParticipant(userID, pollID int64) error {
	_, err := db.mapping.Exec(fmt.Sprintf(
		"delete from %s where %s=$1 and %s=$2;", cParticipantTableName, cUserID,
		cPollID), userID, pollID)
	return err
}
