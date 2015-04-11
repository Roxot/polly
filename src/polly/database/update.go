package database

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

func UpdatePollLastUpdatedTx(pollId int, lastUpdated int64,
	transaction *gorp.Transaction) error {

	_, err := transaction.Exec(fmt.Sprintf("update %s set %s=$1 where %s=%d;",
		cPollTableName, cLastUpdated, cId, pollId), lastUpdated)
	return err
}
