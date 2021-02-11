package database

import (
	"reflect"
	"testing"

	"github.com/roxot/polly"
)

type userDB interface {
	InsertUser(user *polly.User) error
	GetUser(id int64) (*polly.User, error)
	UpdateUser(user *polly.NillableUser) error
	DeleteUser(id int64) error
	CountUsers() (int, error)
}

func TestDBUserCRUD(t *testing.T) {
	if err := testDB.Initialize(); err != nil {
		t.Error("Could not intitialize database:", err)
	}

	testUserCRUD(testDB, t)
	testDB.MustExec(dropTables)
}

func TestTxUserCRUD(t *testing.T) {
	if err := testDB.Initialize(); err != nil {
		t.Error("Could not intitialize database:", err)
	}

	tx, err := testDB.Begin()
	if err != nil {
		t.Error("Could not start transaction", err)
	}

	testUserCRUD(tx, t)

	if err := tx.Commit(); err != nil {
		t.Error("Could not commit transaction", err)
	}

	testDB.MustExec(dropTables)
}

func testUserCRUD(db userDB, t *testing.T) {

	// Test user insertion and retrieval.

	user := polly.User{
		Token:       "SomeSecretToken",
		DisplayName: "Polly User",
		DeviceType:  polly.DEVICE_TYPE_ANDROID,
		DeviceGUID:  "SomeGUID",
		ProfilePic:  "SomeURL",
	}

	if err := db.InsertUser(&user); err != nil {
		t.Error("Failed to insert user:", err)
	}

	savedUser, err := db.GetUser(user.ID)
	if err != nil {
		t.Error("Couldn't find user:", err)
	}

	if !reflect.DeepEqual(user, *savedUser) {
		t.Errorf("Inserted user and saved user not equal: %v, %v",
			user, *savedUser)
	}

	// Test auto incrementing IDs.

	oldID := user.ID
	if err := db.InsertUser(&user); err != nil {
		t.Error("Failed to insert user:", err)
	} else if user.ID <= oldID {
		t.Errorf("ID not incrementing, expected %d, got %d", oldID+1, user.ID)
	}

	// Test updating a user.

	expectedUser := polly.User{
		Token:       "SomeNewSecretToken",
		DisplayName: "SomeNewDisplayName",
		DeviceType:  polly.DEVICE_TYPE_IPHONE,
		DeviceGUID:  "SomeNewGUID",
		ProfilePic:  "SomeNewURL",
	}
	expectedUser.ID = user.ID

	nillableUser := polly.NillableUser{
		Token:       &expectedUser.Token,
		DisplayName: &expectedUser.DisplayName,
		DeviceType:  &expectedUser.DeviceType,
		DeviceGUID:  &expectedUser.DeviceGUID,
		ProfilePic:  &expectedUser.ProfilePic,
	}
	nillableUser.ID = user.ID

	if err := db.UpdateUser(&nillableUser); err != nil {
		t.Error("Failed to update user:", err)
	}

	savedUser, err = db.GetUser(expectedUser.ID)
	if err != nil {
		t.Error("Couldn't find user:", err)
	}

	if !reflect.DeepEqual(expectedUser, *savedUser) {
		t.Errorf("Expected updated user and saved user not equal: %v, %v",
			expectedUser, *savedUser)
	}

	// Test updating a single field of a user.

	expectedUser.ProfilePic = "SomeEvenNewerURL"
	nillableUser = polly.NillableUser{
		ProfilePic: &expectedUser.ProfilePic,
	}
	nillableUser.ID = expectedUser.ID

	if err := db.UpdateUser(&nillableUser); err != nil {
		t.Error("Failed to update user:", err)
	}

	savedUser, err = db.GetUser(expectedUser.ID)
	if err != nil {
		t.Error("Couldn't find user:", err)
	}

	if !reflect.DeepEqual(expectedUser, *savedUser) {
		t.Errorf("Expected updated user and saved user not equal: %v, %v",
			expectedUser, *savedUser)
	}

	// Test counting the number of users.

	if count, err := db.CountUsers(); err != nil {
		t.Error("Failed to count the number of users", err)
	} else if count != 2 {
		t.Errorf("Unexpected user count, expected %d, got %d", 2, count)
	}

	// Test deleting a user.

	if err := db.DeleteUser(user.ID); err != nil {
		t.Error("Failed to delete user from the database", err)
	}

	if count, err := db.CountUsers(); err != nil {
		t.Error("Failed to count the number of users", err)
	} else if count != 1 {
		t.Errorf("Unexpected user count after deletion, expected %d, got %d",
			1, count)
	}
}
