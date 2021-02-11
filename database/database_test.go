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
	DROP TABLE IF EXISTS polls;
	DROP TABLE IF EXISTS users;
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
