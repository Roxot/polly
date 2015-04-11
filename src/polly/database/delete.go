package database

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

func (db *Database) DelVerTokensByPhoneNumber(phoneNo string) error {
	_, err := db.dbMap.Exec(fmt.Sprintf("delete from %s where %s=$1;",
		cVerificationTokensTableName, cPhoneNumber), phoneNo)
	return err
}

func DelVotesForUserTx(userId, pollId int,
	transaction *gorp.Transaction) error {

	_, err := transaction.Exec(fmt.Sprintf(
		"delete from %s where %s=$1 and %s=$2;", cVoteTableName, cUserId,
		cPollId), userId, pollId)
	return err
}
