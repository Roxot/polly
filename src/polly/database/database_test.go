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
	newUsr1 := polly.PrivateUser{}
	newUsr1.PhoneNumber = "0600112233"
	newUsr1.Token = "test_token"
	newUsr1.DisplayName = "Test User"
	newUsr1.DeviceType = 1
	newUsr1.DeviceGUID = "PHONEGUID"
	err := db.AddUser(&newUsr1)
	if err != nil {
		t.Fatalf("Unable to add a user: %s.", err)
	}

	// test whether the id was set correctly
	newUsr2 := polly.PrivateUser{}
	newUsr2.PhoneNumber = "061122334455"
	newUsr2.Token = "test_token_2"
	newUsr2.DisplayName = "Test User 2"
	newUsr2.DeviceType = 0
	newUsr2.DeviceGUID = "PHONEGUID2"
	err = db.AddUser(&newUsr2)
	if newUsr2.Id == 0 {
		t.Fatal("Id not set after adding a user.")
	}

	// test whether the user was inserted using the UserById function
	gotUsr1, err := db.UserById(newUsr1.Id)
	if err != nil {
		t.Fatalf("Failed to find inserted user by id: %s.", err)
	}

	// test whether the inserted user was inserted correctly
	if (gotUsr1.Id != newUsr1.Id) ||
		(gotUsr1.PhoneNumber != newUsr1.PhoneNumber) ||
		(gotUsr1.Token != newUsr1.Token) ||
		(gotUsr1.DisplayName != newUsr1.DisplayName) ||
		(gotUsr1.DeviceType != newUsr1.DeviceType) ||
		(gotUsr1.DeviceGUID != newUsr1.DeviceGUID) {

		t.Fatal("Retrieved user using UserById didn't match inserted user.")
	}

	// test whether the user was inserted using the UserByPhoneNumber function
	gotUsr2, err := db.UserByPhoneNumber(newUsr1.PhoneNumber)
	if err != nil {
		t.Fatal("Failed to find inserted user by phone number.")
	}

	// test whether the inserted user was inserted correctly
	if (gotUsr2.Id != newUsr1.Id) ||
		(gotUsr2.PhoneNumber != newUsr1.PhoneNumber) ||
		(gotUsr2.Token != newUsr1.Token) ||
		(gotUsr2.DisplayName != newUsr1.DisplayName) ||
		(gotUsr2.DeviceType != newUsr1.DeviceType) ||
		(gotUsr2.DeviceGUID != newUsr1.DeviceGUID) {

		t.Fatal("Retrieved user using UserByPhoneNumber didn't match inserted user.")
	}

	// test adding duplicate user
	newUsr3 := polly.PrivateUser{}
	newUsr3.PhoneNumber = "0600112233"
	newUsr3.Token = "test_token_3"
	newUsr3.DisplayName = "Test User 3"
	newUsr3.DeviceType = 1
	newUsr3.DeviceGUID = "PHONEGUID3"
	err = db.AddUser(&newUsr3)
	if err == nil {
		t.Fatal("Was able to add a user with an already existing phone number.")
	}

	// test public user by id
	gotUsr3, err := db.PublicUserById(newUsr2.Id)
	if err != nil {
		t.Fatalf("Unable to retrieve public user by id: %s.", err)
	}

	// test whether they matched
	if (gotUsr3.Id != newUsr2.Id) ||
		(gotUsr3.PhoneNumber != newUsr2.PhoneNumber) ||
		(gotUsr3.DisplayName != newUsr2.DisplayName) {
		t.Fatal("User found with PublicUserById did not match.")
	}

	// test public user by phone number
	gotUsr4, err := db.PublicUserByPhoneNumber(newUsr2.PhoneNumber)
	if err != nil {
		t.Fatalf("Unable to retrieve public user by phone number: %s.", err)
	}

	// test whether they matched
	if (gotUsr4.Id != newUsr2.Id) ||
		(gotUsr4.PhoneNumber != newUsr2.PhoneNumber) ||
		(gotUsr4.DisplayName != newUsr2.DisplayName) {
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
