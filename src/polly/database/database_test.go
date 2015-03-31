package database

import (
	"os"
	"polly/database"
	"testing"
)

var db *database.Database

const (
	cDbUser = "polly"
	cDbPass = "w01V3s"
	cDbName = "pollytestdb"
)

func TestUser(t *testing.T) {

	// test adding a user
	newUser1 := database.PrivateUser{}
	newUser1.PhoneNumber = "0600112233"
	newUser1.Token = "test_token"
	newUser1.DisplayName = "Test User"
	newUser1.DeviceType = 1
	newUser1.DeviceGUID = "PHONEGUID"
	err := db.AddUser(&newUser1)
	if err != nil {
		t.Fatalf("Unable to add a user: %s.", err)
	}

	// test whether the id was set correctly
	newUser2 := database.PrivateUser{}
	newUser2.PhoneNumber = "061122334455"
	newUser2.Token = "test_token_2"
	newUser2.DisplayName = "Test User 2"
	newUser2.DeviceType = 0
	newUser2.DeviceGUID = "PHONEGUID2"
	err = db.AddUser(&newUser2)
	if newUser2.Id == 0 {
		t.Fatal("Id not set after adding a user.")
	}

	// test whether the user was inserted using the UserById function
	gotUser1, err := db.UserById(newUser1.Id)
	if err != nil {
		t.Fatalf("Failed to find inserted user by id: %s.", err)
	}

	// test whether the inserted user was inserted correctly
	if (gotUser1.Id != newUser1.Id) ||
		(gotUser1.PhoneNumber != newUser1.PhoneNumber) ||
		(gotUser1.Token != newUser1.Token) ||
		(gotUser1.DisplayName != newUser1.DisplayName) ||
		(gotUser1.DeviceType != newUser1.DeviceType) ||
		(gotUser1.DeviceGUID != newUser1.DeviceGUID) {

		t.Fatal("Retrieved user using UserById didn't match inserted user.")
	}

	// test whether the user was inserted using the UserByPhoneNumber function
	gotUser2, err := db.UserByPhoneNumber(newUser1.PhoneNumber)
	if err != nil {
		t.Fatal("Failed to find inserted user by phone number.")
	}

	// test whether the inserted user was inserted correctly
	if (gotUser2.Id != newUser1.Id) ||
		(gotUser2.PhoneNumber != newUser1.PhoneNumber) ||
		(gotUser2.Token != newUser1.Token) ||
		(gotUser2.DisplayName != newUser1.DisplayName) ||
		(gotUser2.DeviceType != newUser1.DeviceType) ||
		(gotUser2.DeviceGUID != newUser1.DeviceGUID) {

		t.Fatal("Retrieved user using UserByPhoneNumber didn't match inserted user.")
	}

	// test adding duplicate user
	newUser3 := database.PrivateUser{}
	newUser3.PhoneNumber = "0600112233"
	newUser3.Token = "test_token"
	newUser3.DisplayName = "Test User"
	newUser3.DeviceType = 1
	newUser3.DeviceGUID = "PHONEGUID"
	err = db.AddUser(&newUser3)
	if err == nil {
		t.Fatal("Was able to add a user with an already existing phone number.")
	}

	err = db.DropTablesIfExists()
	if err != nil {
		t.Fatal("Unable to clear database.")
	}
}

func TestMain(m *testing.M) {
	dbConfig := database.DbConfig{}
	dbConfig.DbName = cDbName
	dbConfig.PsqlUserPass = cDbPass
	dbConfig.PsqlUser = cDbUser

	var err error
	db, err = database.New(dbConfig)
	assert(err == nil, "Error creating database")

	err = db.DropTablesIfExists()
	assert(err == nil, "Error clearing database")

	err = db.CreateTablesIfNotExists()
	assert(err == nil, "Error creating tables")

	os.Exit(m.Run())
}

func assert(expr bool, message string) {
	if !expr {
		panic(message)
	}
}
