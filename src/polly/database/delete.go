package database

import "fmt"

func (pollyDb Database) DeleteVerificationTokensByPhoneNumber(
	vt *VerificationToken) error {

	_, err := pollyDb.dbMap.Exec(fmt.Sprintf("delete from %s where %s=$1",
		cVerificationTokensTableName, cPhoneNumber), vt.PhoneNumber)
	return err
}
