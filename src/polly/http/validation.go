package http

import (
	"errors"
	"polly"
	"polly/database"
)

/*
 * Validates a poll message by checking the questions, options and participants.
 * In the case of participants their correct display names and email addresses
 * are set in this function as well.
 */
func isValidPollMessage(db *database.Database, pollMsg *polly.PollMessage,
	creatorID int64) error {

	// validate question type has fitting options
	switch pollMsg.Question.Type {
	case polly.QUESTION_TYPE_MC:
		if pollMsg.Options == nil || len(pollMsg.Options) == 0 {
			return errors.New("Empty options list.")
		}
	case polly.QUESTION_TYPE_OPEN:
		if pollMsg.Options != nil && len(pollMsg.Options) > 0 {
			return errors.New("Non-empty option list in open question.")
		}
	default:
		return errors.New("Unknown poll type.")
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
	participantsMap := make(map[int64]bool)
	numParticipants := len(pollMsg.Participants)
	for i := 0; i < numParticipants; i++ {

		// check for duplicate participants
		_, ok := participantsMap[pollMsg.Participants[i].ID]
		if ok {
			return errors.New("Duplicate participant.")
		}

		// check if user exists
		dbUser, err := db.GetUserByID(pollMsg.Participants[i].ID)
		if err != nil {
			return errors.New("Unknown participant.")
		} else {
			pollMsg.Participants[i].DisplayName = dbUser.DisplayName
		}

		// check if user is creator
		if pollMsg.Participants[i].ID == creatorID {
			containsCreator = true
		}

		// add participant to map of particpants
		participantsMap[pollMsg.Participants[i].ID] = true
	}

	// make sure user is a participant
	if !containsCreator {
		return errors.New("Creator not in participants list.")
	}

	return nil
}

func isValidDeviceType(deviceType int) bool {
	return (deviceType == polly.DEVICE_TYPE_ANDROID ||
		deviceType == polly.DEVICE_TYPE_IPHONE)
}
