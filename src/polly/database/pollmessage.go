package database

import (
	"polly"
)

func (db *Database) InsertPollMessage(pollMsg *polly.PollMessage) error {
	var err error

	// start the transaction
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert the poll object
	err = AddPollTX(&pollMsg.MetaData, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update the poll message
	pollMsg.Question.PollID = pollMsg.MetaData.ID

	// insert the question
	err = AddQuestionTX(&pollMsg.Question, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert the options
	numOptions := len(pollMsg.Options)
	for i := 0; i < numOptions; i++ {
		option := &(pollMsg.Options[i])
		option.QuestionID = pollMsg.Question.ID
		option.PollID = pollMsg.MetaData.ID
		err = AddOptionTX(option, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// insert the participants
	numParticipants := len(pollMsg.Participants)
	for i := 0; i < numParticipants; i++ {
		user := pollMsg.Participants[i]
		participant := polly.Participant{} // TODO style inconsistency
		participant.UserID = user.ID
		participant.PollID = pollMsg.MetaData.ID
		err = AddParticipantTX(&participant, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ConstructPollMessage(pollID int64) (*polly.PollMessage,
	error) {

	pollMsg := polly.PollMessage{}

	// retrieve the poll object
	poll, err := db.GetPollByID(pollID)
	pollMsg.MetaData = *poll
	if err != nil {
		return nil, err
	}

	// retrieve the questions
	question, err := db.GetQuestionByPollID(pollID)
	pollMsg.Question = *question
	if err != nil {
		return nil, err
	}

	// retrieve the options
	options, err := db.GetOptionsByPollID(pollID)
	pollMsg.Options = options
	if err != nil {
		return nil, err
	}

	// retrieve the votes
	votes, err := db.GetVotesByPollID(pollID)
	pollMsg.Votes = votes
	if err != nil {
		return nil, err
	}

	// retrieve the participants
	participants, err := db.GetParticipantsByPollID(pollID)
	if err != nil {
		return nil, err
	}

	// convert the participants to user objects
	numParticipants := len(participants)
	pollMsg.Participants = make([]polly.PublicUser, numParticipants)
	var user *polly.PublicUser
	for i := 0; i < numParticipants; i++ {
		user, err = db.GetPublicUserByID(participants[i].UserID)
		if err != nil {
			return nil, err
		}

		pollMsg.Participants[i] = *user
	}

	return &pollMsg, nil
}
