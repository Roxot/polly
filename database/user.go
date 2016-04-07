package database

import "github.com/roxot/polly"

func (db *DB) GetUser(id int64) (*polly.User, error) {
	user := new(polly.User)
	err := db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func (db *DB) UpdateUser(user *polly.NillableUser) error {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	if user.Token != nil {
		_, err = tx.Exec("UPDATE users SET token=$1 WHERE id=$2",
			*(user.Token), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.DisplayName != nil {
		_, err = tx.Exec("UPDATE users SET display_name=$1 WHERE id=$2",
			*(user.DisplayName), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.DeviceType != nil {
		_, err = tx.Exec("UPDATE users SET device_type=$1 WHERE id=$2",
			*(user.DeviceType), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.DeviceGUID != nil {
		_, err = tx.Exec("UPDATE users SET device_guid=$1 WHERE id=$2",
			*(user.DeviceGUID), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.ProfilePic != nil {
		_, err = tx.Exec("UPDATE users SET profile_pic=$1 WHERE id=$2",
			*(user.ProfilePic), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (db *DB) InsertUser(user *polly.User) error {
	return db.QueryRow("INSERT INTO users (token, display_name, device_type, device_guid, profile_pic) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Token, user.DisplayName, user.DeviceType, user.DeviceGUID,
		user.ProfilePic).Scan(&user.ID)
}

func (db *DB) DeleteUser(id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

func (db *DB) CountUsers() (int, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) FROM users")
	return count, err
}
