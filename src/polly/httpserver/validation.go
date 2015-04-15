package httpserver

import (
	"errors"
	"polly"
	"polly/database"
	"unicode"
)

func isValidPhoneNumber(phoneNo string) bool {
	if len(phoneNo) != 10 {
		return false
	}

	for idx, val := range phoneNo {
		if idx == 0 {
			if val != '0' {
				return false
			}
		} else if idx == 1 {
			if val != '6' {
				return false
			}
		} else if !unicode.IsNumber(val) {
			return false
		}
	}

	return true
}

/*
 * Validates a poll message by checking the questions, options and participants.
 * In the case of participants their correct display names and phone numbers are
 * set in this function as well.
 */
func isValidPollMessage(db *database.Database, pollMsg *PollMessage,
	creatorId int) error {

	// validate question type has fitting options
	switch pollMsg.Question.Type {
	case polly.QUESTION_TYPE_MC:
		if pollMsg.Options == nil || len(pollMsg.Options) == 0 {
			return errors.New("Empty options list.")
		}
	case polly.QUESTION_TYPE_OP:
		if pollMsg.Options != nil && len(pollMsg.Options) > 0 {
			return errors.New("Non-empty option list in open question.")
		}
	case polly.QUESTION_TYPE_DT:
		// TODO no support yet for date polls
	}

	// don't accept empty question titles
	if len(pollMsg.Question.Title) == 0 {
		return errors.New("Empty question title.")
	}

	// don't accept empty option values
	numOptions := len(pollMsg.Options)
	for i := 0; i < numOptions; i++ {
		if len(pollMsg.Options[i].Value) == 0 {
			return errors.New("Empty value field in option object.")
		}
	}

	containsCreator := false
	particMap := make(map[int]bool)
	numPartics := len(pollMsg.Participants)
	for i := 0; i < numPartics; i++ {

		// check for duplicate participants
		_, ok := particMap[pollMsg.Participants[i].Id]
		if ok {
			return errors.New("Duplicate participant.")
		}

		// check if user exists
		dbUser, err := db.UserById(pollMsg.Participants[i].Id)
		if err != nil {
			return errors.New("Unknown participant.")
		} else {
			pollMsg.Participants[i].DisplayName = dbUser.DisplayName
			pollMsg.Participants[i].PhoneNumber = dbUser.PhoneNumber
		}

		// check if user is creator
		if pollMsg.Participants[i].Id == creatorId {
			containsCreator = true
		}

		// add participant to map of particpants
		particMap[pollMsg.Participants[i].Id] = true
	}

	// make sure user is a participant
	if !containsCreator {
		return errors.New("Creator not in participants list.")
	}

	return nil
}
