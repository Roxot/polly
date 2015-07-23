package http

import (
	"polly"
	"polly/database"
	"time"
)

type PollMessage struct {
	MetaData     polly.Poll         `json:"meta_data"`
	Question     polly.Question     `json:"question"`
	Options      []polly.Option     `json:"options"`
	Votes        []polly.Vote       `json:"votes"`
	Participants []polly.PublicUser `json:"participants"`
}

func (server *sServer) InsertPollMessage(pollMsg *PollMessage) error {
	var err error

	// start the transaction
	transaction, err := server.db.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}

	// set the time creation date and last update time to now
	now := time.Now().Unix()
	pollMsg.MetaData.CreationDate = now
	pollMsg.MetaData.LastUpdated = now

	// insert the poll object
	err = database.AddPollTX(&pollMsg.MetaData, transaction)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// update the poll message
	pollMsg.Question.PollID = pollMsg.MetaData.ID

	// insert the question
	err = database.AddQuestionTX(&pollMsg.Question, transaction)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// insert the options
	numOptions := len(pollMsg.Options)
	for i := 0; i < numOptions; i++ {
		option := &(pollMsg.Options[i])
		option.QuestionID = pollMsg.Question.ID
		option.PollID = pollMsg.MetaData.ID
		err = database.AddOptionTX(option, transaction)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	// insert the participants
	numParticipants := len(pollMsg.Participants)
	for i := 0; i < numParticipants; i++ {
		user := pollMsg.Participants[i]
		partic := polly.Participant{}
		partic.UserID = user.ID
		partic.PollID = pollMsg.MetaData.ID
		err = database.AddParticipantTX(&partic, transaction)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (server *sServer) ConstructPollMessage(pollID int) (*PollMessage, error) {
	pollMsg := PollMessage{}

	// retrieve the poll object
	poll, err := server.db.GetPollByID(pollID)
	pollMsg.MetaData = *poll
	if err != nil {
		return nil, err
	}

	// retrieve the questions
	question, err := server.db.GetQuestionByPollID(pollID)
	pollMsg.Question = *question
	if err != nil {
		return nil, err
	}

	// retrieve the options
	options, err := server.db.GetOptionsByPollID(pollID)
	pollMsg.Options = options
	if err != nil {
		return nil, err
	}

	// retrieve the votes
	votes, err := server.db.GetVotesByPollID(pollID)
	pollMsg.Votes = votes
	if err != nil {
		return nil, err
	}

	// retrieve the participants
	participants, err := server.db.GetParticipantsByPollID(pollID)
	if err != nil {
		return nil, err
	}

	// convert the participants to user objects
	numParticipants := len(participants)
	pollMsg.Participants = make([]polly.PublicUser, numParticipants)
	var user *polly.PublicUser
	for i := 0; i < numParticipants; i++ {
		user, err = server.db.GetPublicUserByID(participants[i].UserID)
		if err != nil {
			return nil, err
		}

		pollMsg.Participants[i] = *user
	}

	return &pollMsg, nil
}
