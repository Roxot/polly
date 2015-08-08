package http

import (
	"polly"
	"polly/database"
	"strings"
)

/*
 * Validates a poll message by checking the questions, options and participants.
 * In the case of participants their correct display names and email addresses
 * are set in this function as well.
 */
func isValidPollMessage(db *database.Database, pollMsg *polly.PollMessage,
	creatorID int64) int {

	// validate question type has fitting options
	switch pollMsg.Question.Type {
	case polly.QUESTION_TYPE_MC:
		if pollMsg.Options == nil || len(pollMsg.Options) == 0 {
			return ERR_BAD_EMPTY_POLL
		}
	case polly.QUESTION_TYPE_OPEN:
		// skip
	default:
		return ERR_BAD_POLL_TYPE
	}

	// don't accept empty question titles
	pollMsg.Question.Title = strings.TrimSpace(pollMsg.Question.Title)
	if len(pollMsg.Question.Title) == 0 {
		return ERR_BAD_EMPTY_QUESTION
	}

	// don't accept empty option values
	numOptions := len(pollMsg.Options)
	for i := 0; i < numOptions; i++ {
		pollMsg.Options[i].Value = strings.TrimSpace(pollMsg.Options[i].Value)
		if len(pollMsg.Options[i].Value) == 0 {
			return ERR_BAD_EMPTY_OPTION
		}
	}

	containsCreator := false
	participantsMap := make(map[int64]bool)
	numParticipants := len(pollMsg.Participants)
	for i := 0; i < numParticipants; i++ {

		// check for duplicate participants
		_, ok := participantsMap[pollMsg.Participants[i].ID]
		if ok {
			return ERR_BAD_DUPLICATE_PARTICIPANT
		}

		// check if user exists
		dbUser, err := db.GetUserByID(pollMsg.Participants[i].ID)
		if err != nil {
			return ERR_BAD_NO_USER
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
		return ERR_BAD_NO_CREATOR
	}

	return NO_ERR
}

func isValidDeviceType(deviceType int) bool {
	return (deviceType == polly.DEVICE_TYPE_ANDROID ||
		deviceType == polly.DEVICE_TYPE_IPHONE)
}
