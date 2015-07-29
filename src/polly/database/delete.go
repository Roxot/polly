package database

import (
	"fmt"

	"polly/internal/gopkg.in/gorp.v1"
)

func DeleteVotesForUserTX(userID, pollID int64, tx *gorp.Transaction) error {
	_, err := tx.Exec(fmt.Sprintf("delete from %s where %s=$1 and %s=$2;",
		cVoteTableName, cUserID, cPollID), userID, pollID)
	return err
}
