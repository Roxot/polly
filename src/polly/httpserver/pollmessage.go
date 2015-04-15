package httpserver

import (
	"polly"
	"time"
)

type PollMessage struct {
	MetaData     polly.Poll         `json:"meta_data"`
	Question     polly.Question     `json:"question"`
	Options      []polly.Option     `json:"options"`
	Votes        []polly.Vote       `json:"votes"`
	Participants []polly.PublicUser `json:"participants"`
}

func (srv *HTTPServer) InsertPollMessage(pollMsg *PollMessage) error {
	var err error

	// start the transaction
	transaction, err := srv.db.Begin()
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		} else {
			return err
		}
	}

	// set the time creation date and last update time to now
	now := time.Now().Unix()
	pollMsg.MetaData.CreationDate = now
	pollMsg.MetaData.LastUpdated = now

	// insert the poll object
	err = srv.db.AddPollTx(&pollMsg.MetaData, transaction)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		} else {
			return err
		}
	}

	// update the poll message
	pollMsg.Question.PollId = pollMsg.MetaData.Id

	// insert the question
	err = srv.db.AddQuestionTx(&pollMsg.Question, transaction)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		} else {
			return err
		}
	}

	// insert the options
	numOptions := len(pollMsg.Options)
	for i := 0; i < numOptions; i++ {
		option := &(pollMsg.Options[i])
		option.QuestionId = pollMsg.Question.Id
		option.PollId = pollMsg.MetaData.Id
		err = srv.db.AddOptionTx(option, transaction)
		if err != nil {
			rollbackErr := transaction.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			} else {
				return err
			}
		}
	}

	// insert the participants
	numPartics := len(pollMsg.Participants)
	for i := 0; i < numPartics; i++ {
		usr := pollMsg.Participants[i]
		partic := polly.Participant{}
		partic.UserId = usr.Id
		partic.PollId = pollMsg.MetaData.Id
		err = srv.db.AddParticipantTx(&partic, transaction)
		if err != nil {
			rollbackErr := transaction.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			} else {
				return err
			}
		}
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (srv *HTTPServer) ConstructPollMessage(pollId int) (*PollMessage, error) {
	pollMsg := PollMessage{}

	// retrieve the poll object
	poll, err := srv.db.PollById(pollId)
	pollMsg.MetaData = *poll
	if err != nil {
		return nil, err
	}

	// retrieve the questions
	question, err := srv.db.QuestionByPollId(pollId)
	pollMsg.Question = *question
	if err != nil {
		return nil, err
	}

	// retrieve the options
	options, err := srv.db.OptionsByPollId(pollId)
	pollMsg.Options = options
	if err != nil {
		return nil, err
	}

	// retrieve the votes
	votes, err := srv.db.VotesByPollId(pollId)
	pollMsg.Votes = votes
	if err != nil {
		return nil, err
	}

	// retrieve the participants
	participants, err := srv.db.ParticipantsByPollId(pollId)
	if err != nil {
		return nil, err
	}

	// convert the participants to user objects
	numPartics := len(participants)
	pollMsg.Participants = make([]polly.PublicUser, numPartics)
	var user *polly.PublicUser
	for i := 0; i < numPartics; i++ {
		user, err = srv.db.PublicUserById(participants[i].UserId)
		if err != nil {
			return nil, err
		}

		pollMsg.Participants[i] = *user
	}

	return &pollMsg, nil
}
