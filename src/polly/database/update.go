package database

import (
	"fmt"

	"polly/internal/gopkg.in/gorp.v1"
)

func UpdatePollLastUpdatedTX(pollID, lastUpdated int64,
	tx *gorp.Transaction) error {

	_, err := tx.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cPollTableName, cLastUpdated, cID), lastUpdated, pollID)
	return err
}

func (db *Database) UpdateUser(userID int64, displayName,
	deviceGUID string) error {

	_, err := db.mapping.Exec(fmt.Sprintf(
		"update %s set %s=$1, %s=$2 where %s=$3;", cUserTableName,
		cDisplayName, cDeviceGUID, cID), displayName, deviceGUID, userID)
	return err
}

func (db *Database) UpdateDisplayName(userID int64, displayName string) error {
	_, err := db.mapping.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cUserTableName, cDisplayName, cID), displayName, userID)
	return err
}

func (db *Database) UpdateDeviceGUID(userID int64, deviceGUID string) error {
	_, err := db.mapping.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cUserTableName, cDeviceGUID, cID), deviceGUID, userID)
	return err
}

func (db *Database) UpdateToken(userID int64, token string) error {
	_, err := db.mapping.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cUserTableName, cToken, cID), token, userID)
	return err
}

func (db *Database) UpdateSequenceNumber(pollID int64) error {
	_, err := db.mapping.Exec(fmt.Sprintf("update %d set %s=%s+1 where %s=$1;",
		cPollTableName, cSequenceNumber, cSequenceNumber, cID), pollID)
	return err
}

/* Updates a poll its last updated and sequence number. */
func UpdatePollTX(pollID, lastUpdated int64, tx *gorp.Transaction) error {
	_, err := tx.Exec(fmt.Sprintf(
		"update %s set %s=%s+1, %s=$1 where %s=$2;", cPollTableName,
		cSequenceNumber, cSequenceNumber, cLastUpdated, cID), lastUpdated,
		pollID)
	return err
}

func UpdateOptionSequenceNumberTX(optionID int64, sequenceNumber int,
	tx *gorp.Transaction) error {
	_, err := tx.Exec(fmt.Sprintf("update %s set %s=$1 where %s=$2;",
		cOptionTableName, cSequenceNumber, cID),
		sequenceNumber, optionID)
	return err
}

func UpdateSequenceNumberTX(pollID int64, tx *gorp.Transaction) error {
	_, err := tx.Exec(fmt.Sprintf("update %s set %s=%s+1 where %s=$1;",
		cPollTableName, cSequenceNumber, cSequenceNumber, cID), pollID)
	return err
}
