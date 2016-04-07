package database

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/roxot/polly"
)

func TestUserCRUD(t *testing.T) {
	err := testDB.Initialize()
	if err != nil {
		t.Error("Could not intitialize database:", err)
	}

	// Test user insertion and retrieval.

	user := polly.User{
		Token:       "SomeSecretToken",
		DisplayName: "Polly User",
		DeviceType:  polly.DEVICE_TYPE_ANDROID,
		DeviceGUID:  "SomeGUID",
		ProfilePic:  "SomeURL",
	}

	err = testDB.InsertUser(&user)
	if err != nil {
		t.Error("Failed to insert user:", err)
	}

	savedUser, err := testDB.GetUser(user.ID)
	if err != nil {
		t.Error("Couldn't find user:", err)
	}

	if !reflect.DeepEqual(user, *savedUser) {
		t.Errorf("Inserted user and saved user not equal: %v, %v",
			user, *savedUser)
	}

	// Test auto incrementing IDs.

	oldID := user.ID
	err = testDB.InsertUser(&user)
	if err != nil {
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

	err = testDB.UpdateUser(&nillableUser)
	if err != nil {
		t.Error("Failed to update user:", err)
	}

	savedUser, err = testDB.GetUser(expectedUser.ID)
	if err != nil {
		t.Error("Couldn't find user:", err)
	}

	if !reflect.DeepEqual(expectedUser, *savedUser) {
		t.Errorf("Inserted user and saved user not equal: %v, %v",
			expectedUser, *savedUser)
	}

	// Test updating a single field of a user.

	expectedUser.ProfilePic = "SomeEvenNewerURL"
	nillableUser = polly.NillableUser{
		ProfilePic: &expectedUser.ProfilePic,
	}

	fmt.Println(nillableUser, *(nillableUser.ProfilePic))
	err = testDB.UpdateUser(&nillableUser)
	if err != nil {
		t.Error("Failed to update user:", err)
	}

	savedUser, err = testDB.GetUser(expectedUser.ID)
	if err != nil {
		t.Error("Couldn't find user:", err)
	}

	if !reflect.DeepEqual(expectedUser, *savedUser) {
		t.Errorf("Inserted user and saved user not equal: %v, %v",
			expectedUser, *savedUser)
	}
}
