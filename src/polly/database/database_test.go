package database

import (
	"os"
	"polly"
	"polly/database"
	"testing"
)

var db *database.Database

const (
	cDbUser = "polly"
	cDbPass = "w01V3s"
	cDbName = "pollytestdb"
)

func TestUsers(t *testing.T) {

	// test adding a user
	newUser1 := polly.PrivateUser{}
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
	newUser2 := polly.PrivateUser{}
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
	newUser3 := polly.PrivateUser{}
	newUser3.PhoneNumber = "0600112233"
	newUser3.Token = "test_token_3"
	newUser3.DisplayName = "Test User 3"
	newUser3.DeviceType = 1
	newUser3.DeviceGUID = "PHONEGUID3"
	err = db.AddUser(&newUser3)
	if err == nil {
		t.Fatal("Was able to add a user with an already existing phone number.")
	}

	// test public user by id
	gotUser3, err := db.PublicUserById(newUser2.Id)
	if err != nil {
		t.Fatalf("Unable to retrieve public user by id: %s.", err)
	}

	// test whether they matched
	if (gotUser3.Id != newUser2.Id) ||
		(gotUser3.PhoneNumber != newUser2.PhoneNumber) ||
		(gotUser3.DisplayName != newUser2.DisplayName) {
		t.Fatal("User found with PublicUserById did not match.")
	}

	// test public user by phone number
	gotUser4, err := db.PublicUserByPhoneNumber(newUser2.PhoneNumber)
	if err != nil {
		t.Fatalf("Unable to retrieve public user by phone number: %s.", err)
	}

	// test whether they matched
	if (gotUser4.Id != newUser2.Id) ||
		(gotUser4.PhoneNumber != newUser2.PhoneNumber) ||
		(gotUser4.DisplayName != newUser2.DisplayName) {
		t.Fatal("User found with PublicUserById did not match.")
	}

	// clear the database
	err = db.DropTablesIfExists()
	if err != nil {
		t.Fatal("Unable to clear database.")
	}
}

func TestPolls(t *testing.T) {
	t.SkipNow()
}

func TestQuestions(t *testing.T) {
	t.SkipNow()
}

func TestOptions(t *testing.T) {
	t.SkipNow()
}

func TestVotes(t *testing.T) {
	t.SkipNow()
}

func TestParticipants(t *testing.T) {
	t.SkipNow()
}

func TestVerTokens(t *testing.T) {
	t.SkipNow()
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
