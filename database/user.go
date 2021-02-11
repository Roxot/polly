package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/roxot/polly"
)

// InsertUser using this DB.
func (db *DB) InsertUser(user *polly.User) error {
	return InsertUser(db, user)
}

// GetUser using this DB.
func (db *DB) GetUser(id int64) (*polly.User, error) {
	return GetUser(db, id)
}

// UpdateUser using this DB.
func (db *DB) UpdateUser(user *polly.NillableUser) error {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := UpdateUser(db, user); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// DeleteUser using this DB.
func (db *DB) DeleteUser(id int64) error {
	return DeleteUser(db, id)
}

// CountUsers using this DB.
func (db *DB) CountUsers() (int, error) {
	return CountUsers(db)
}

// InsertUser within a transaction.
func (tx *Tx) InsertUser(user *polly.User) error {
	return InsertUser(tx, user)
}

// GetUser within a transaction.
func (tx *Tx) GetUser(id int64) (*polly.User, error) {
	return GetUser(tx, id)
}

// UpdateUser within a transaction.
func (tx *Tx) UpdateUser(user *polly.NillableUser) error {
	return UpdateUser(tx, user)
}

// DeleteUser within a transaction.
func (tx *Tx) DeleteUser(id int64) error {
	return DeleteUser(tx, id)
}

// CountUsers within a transaction.
func (tx *Tx) CountUsers() (int, error) {
	return CountUsers(tx)
}

// InsertUser inserts the user using the provided sqlx.Queryer and sets the ID
// field on user.
func InsertUser(q sqlx.Queryer, user *polly.User) error {
	return q.QueryRowx("INSERT INTO users (token, display_name, device_type, device_guid, profile_pic) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Token, user.DisplayName, user.DeviceType, user.DeviceGUID,
		user.ProfilePic).Scan(&user.ID)
}

// GetUser returns the user with the provided id using the provided
// sqlx.Execer.
func GetUser(q sqlx.Queryer, id int64) (*polly.User, error) {
	user := new(polly.User)
	err := sqlx.Get(q, user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

// UpdateUser updates all non-nil fields in user for the user with id user.ID
// using the provided sqlx.Execer.
func UpdateUser(e sqlx.Execer, user *polly.NillableUser) error {
	if user.Token != nil {
		if _, err := e.Exec("UPDATE users SET token=$1 WHERE id=$2",
			*(user.Token), user.ID); err != nil {
			return err
		}
	}

	if user.DisplayName != nil {
		if _, err := e.Exec("UPDATE users SET display_name=$1 WHERE id=$2",
			*(user.DisplayName), user.ID); err != nil {
			return err
		}
	}

	if user.DeviceType != nil {
		if _, err := e.Exec("UPDATE users SET device_type=$1 WHERE id=$2",
			*(user.DeviceType), user.ID); err != nil {
			return err
		}
	}

	if user.DeviceGUID != nil {
		if _, err := e.Exec("UPDATE users SET device_guid=$1 WHERE id=$2",
			*(user.DeviceGUID), user.ID); err != nil {
			return err
		}
	}

	if user.ProfilePic != nil {
		if _, err := e.Exec("UPDATE users SET profile_pic=$1 WHERE id=$2",
			*(user.ProfilePic), user.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteUser deletes a user with the provided id using the provided
// sqlx.Execer.
func DeleteUser(e sqlx.Execer, id int64) error {
	_, err := e.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

// CountUsers returns the number of users in the database using the provided
// sqlx.Queryer.
func CountUsers(q sqlx.Queryer) (int, error) {
	var count int
	err := sqlx.Get(q, &count, "SELECT count(*) FROM users")
	return count, err
}
