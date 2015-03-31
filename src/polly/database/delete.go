package database

import "fmt"

func (db *Database) DelVerTokensByPhoneNumber(vt *VerToken) error {
	_, err := db.dbMap.Exec(fmt.Sprintf("delete from %s where %s=$1",
		cVerificationTokensTableName, cPhoneNumber), vt.PhoneNumber)
	return err
}
