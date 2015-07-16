package database

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

func UpdatePollLastUpdatedTx(pollID int, lastUpdated int64,
	tx *gorp.Transaction) error {

	_, err := tx.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cPollTableName, cLastUpdated, cID), lastUpdated, pollID)
	return err
}

func (db *Database) UpdateUser(userID int, displayName,
	deviceGUID string) error {

	_, err := db.mapping.Exec(fmt.Sprintf(
		"update %s set %s=$1, %s=$2 where %s=$3;", cUserTableName,
		cDisplayName, cDeviceGUID, cId), dspName, deviceGUID, userID)
	return err
}

func (db *Database) UpdateDisplayName(userID int, displayName string) error {
	_, err := db.mapping.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cUserTableName, cDisplayName, cId), dspName, userID)
	return err
}

func (db *Database) UpdateDeviceGUID(userID int, deviceGUID string) error {
	_, err := db.mapping.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cUserTableName, cDeviceGUID, cId), deviceGUID, userID)
	return err
}
