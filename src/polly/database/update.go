package database

import (
	"fmt"
	"log"

	"gopkg.in/gorp.v1"
)

func UpdatePollLastUpdatedTx(pollId int, lastUpdated int64,
	transaction *gorp.Transaction) error {

	_, err := transaction.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cPollTableName, cLastUpdated, cId), pollId, lastUpdated)
	return err
}

func (db *Database) UpdateUser(usrID int, dspName, dvcGUID string) error {
	_, err := db.dbMap.Exec(fmt.Sprintf(
		"update %s set %s=$1, %s=$2 where %s=$3;", cUserTableName,
		cDisplayName, cDeviceGUID, cId), dspName, dvcGUID, usrID)
	log.Println("Updating to", dvcGUID)
	return err
}
