package database

import "fmt"

type PollMetaData struct {
	PollId       int    `json:"poll_id, omitempty"`
	CreationDate int64  `json:"creation_date, omitempty"`
	Title        string `json:"title"`
}

type PollData struct {
	MetaData     PollMetaData `json:"meta_data"`
	Questions    []Question   `json:"questions"`
	Options      []Option     `json:"options"`
	Votes        []Vote       `json:"votes, omitempty"`
	Creator      User         `json:"creator, omitempty"`
	Participants []User       `json:"participants, omitempty"`
}

func (pollyDb Database) RetrievePollData(pollId int) (PollData, error) {
	pollData := PollData{}

	poll, err := pollyDb.FindPollById(pollId)
	if err != nil {
		return pollData, nil
	}

	pollData.MetaData = PollMetaData{poll.Id, poll.CreationDate, poll.Title}

	questions, err := pollyDb.FindQuestionsByPollId(pollId)
	pollData.Questions = questions
	if err != nil {
		return pollData, err
	}

	options, err := pollyDb.FindOptionsByPollId(pollId)
	pollData.Options = options
	if err != nil {
		return pollData, err
	}

	votes, err := pollyDb.FindVotesByPollId(pollId)
	pollData.Votes = votes
	if err != nil {
		return pollData, err
	}

	creator, err := pollyDb.FindUserById(poll.CreatorId)
	pollData.Creator = creator
	if err != nil {
		return pollData, err
	}

	participants, err := pollyDb.FindParticipantsByPollId(pollId)
	if err != nil {
		return pollData, err
	}

	pollData.Participants = make([]User, len(participants))
	var user User
	for index, participant := range participants {
		user, err = pollyDb.FindUserById(participant.UserId)
		if err != nil {
			return pollData, err
		}

		pollData.Participants[index] = user
	}

	return pollData, nil
}

// TODO rollback if fails
// Ignore votes
// rollback err highest priority
func (pollyDb Database) InsertPollData(pollData *PollData) (error, bool) {
	var err error

	transaction, err := pollyDb.dbMap.Begin()
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return rollbackErr, true
		} else {
			return err, true
		}
	}

	poll := Poll{
		CreatorId:    pollData.Creator.Id,
		CreationDate: pollData.MetaData.CreationDate,
		Title:        pollData.MetaData.Title}
	err = transaction.Insert(&poll)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return rollbackErr, true
		} else {
			return err, true
		}
	}

	// Set server-side poll id
	pollData.MetaData.PollId = poll.Id

	for index, question := range pollData.Questions {
		pollData.Questions[index].PollId = poll.Id
		pollData.Questions[index].ClientId = question.Id
		err = transaction.Insert(&question)
		if err != nil {
			rollbackErr := transaction.Rollback()
			if rollbackErr != nil {
				return rollbackErr, true
			} else {
				return err, true
			}
		}

		// update question id
		pollData.Questions[index].Id = question.Id
	}

OuterLoop:
	for index, option := range pollData.Options {
		for _, question := range pollData.Questions {
			if option.QuestionId == question.ClientId {
				pollData.Options[index].PollId = poll.Id
				pollData.Options[index].QuestionId = question.Id
				err = transaction.Insert(&option)
				if err != nil {
					rollbackErr := transaction.Rollback()
					if rollbackErr != nil {
						return rollbackErr, true
					} else {
						return err, true
					}
				}

				pollData.Options[index].Id = option.Id
				continue OuterLoop
			}
		}

		// no matching question id found
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return rollbackErr, true
		} else {
			return fmt.Errorf("No question found with id %d", option.Id), false
		}
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		return err, true
	}

	// TODO participants
	// ...

	return nil, false
}
