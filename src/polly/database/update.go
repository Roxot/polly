package database

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

func UpdatePollLastUpdatedTx(pollId int, lastUpdated int64,
	transaction *gorp.Transaction) error {

	_, err := transaction.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cPollTableName, cLastUpdated, cId), lastUpdated, pollId)
	return err
}

func (db *Database) UpdateUser(usrID int, dspName, dvcGUID string) error {
	_, err := db.dbMap.Exec(fmt.Sprintf(
		"update %s set %s=$1, %s=$2 where %s=$3;", cUserTableName,
		cDisplayName, cDeviceGUID, cId), dspName, dvcGUID, usrID)
	return err
}

func (db *Database) UpdateDisplayName(usrID int, dspName string) error {
	_, err := db.dbMap.Exec(fmt.Sprintf(
		"update %s set %s=$1 where %s=$2;", cUserTableName,
		cDisplayName, cId), dspName, usrID)
	return err
}

func (db *Database) UpdateDeviceGUID(usrID int, dvcGUID string) error {
	_, err := db.dbMap.Exec(fmt.Sprintf(
		"update %s set %s=$1 where %s=$2;", cUserTableName,
		cDeviceGUID, cId), dvcGUID, usrID)
	return err
}
