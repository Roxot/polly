package database

import "fmt"

func (db *Database) DelVerTokensByPhoneNumber(phoneNo string) error {
	_, err := db.dbMap.Exec(fmt.Sprintf("delete from %s where %s=$1",
		cVerificationTokensTableName, cPhoneNumber), phoneNo)
	return err
}
