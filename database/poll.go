package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/roxot/polly"
)

func InsertUsers(q sqlx.Queryer, poll *polly.Poll) error {
	return q.QueryRowx("INSERT INTO polls (creator_id, creation_date, closing_date, last_updated, sequence_number, last_event_user, last_event_user_id, last_event_title, last_event_type) VALUES ($1, now(), $2, now(), 0, $3, $4, $5, $%) RETURNING id, creation_date, last_updated, sequence_number",
		poll.CreatorID, poll.ClosingDate, poll.LastEventUser,
		poll.LastEventUserID, poll.LastEventTitle, poll.LastEventType).
		Scan(&poll.ID, &poll.CreationDate, &poll.LastUpdated,
			poll.SequenceNumber)
}
