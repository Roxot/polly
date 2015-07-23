package database

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

func (db *Database) DeleteVerTokensByEmail(email string) error {
	_, err := db.mapping.Exec(fmt.Sprintf("delete from %s where %s=$1;",
		cVerificationTokensTableName, cEmail), email)
	return err
}

func DeleteVotesForUserTx(userID, pollID int, tx *gorp.Transaction) error {
	_, err := tx.Exec(fmt.Sprintf("delete from %s where %s=$1 and %s=$2;",
		cVoteTableName, cUserID, cPollID), userID, pollID)
	return err
}
