package database

import (
	"flag"
	"log"
	"os"
	"testing"
)

var testConfig = &Config{
	User:     "pollytest",
	Password: "",
	DBName:   "pollytestdb",
	SSLMode:  "disable",
}

var dropTables = `
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS participants;
	DROP TABLE IF EXISTS polls;
	DROP TABLE IF EXISTS questions;
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS verification_tokens;
	DROP TABLE IF EXISTS votes;
`

var testDB *DB

func TestMain(m *testing.M) {
	var err error
	flag.Parse()

	testDB, err = Connect(testConfig)
	if err != nil {
		log.Fatal(err)
	}
	testDB.MustExec(dropTables)

	os.Exit(m.Run())
}
