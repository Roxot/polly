package http

import (
	"encoding/json"
	"net/http"
	"polly"
	"polly/database"
	"time"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cVoteTag = "POST/VOTE"
)

func (server *sServer) Vote(writer http.ResponseWriter, request *http.Request,
	_ httprouter.Params) {
	var err error

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cVoteTag, writer, request)
		return
	}

	// decode the vote message
	var voteMsg polly.VoteMessage
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&voteMsg)
	if err != nil {
		server.respondWithError(ERR_BAD_JSON, err, cVoteTag, writer, request)
		return
	}

	// retrieve the poll id belonging to the option or question id
	var pollID int64
	switch voteMsg.Type {
	case polly.VOTE_TYPE_NEW:
		question, err := server.db.GetQuestionByID(voteMsg.ID)
		if err != nil {
			server.respondWithError(ERR_BAD_NO_QUESTION, err, cVoteTag, writer,
				request)
			return

		}

		pollID = question.ID
		if question.Type != polly.QUESTION_TYPE_OPEN {
			server.respondWithError(ERR_ILL_ADD_OPTION, nil, cVoteTag, writer,
				request)
			return
		} else if len(voteMsg.Value) == 0 {
			server.respondWithError(ERR_BAD_EMPTY_OPTION, nil, cVoteTag, writer,
				request)
			return
		}

	case polly.VOTE_TYPE_UPVOTE:
		pollID, err = server.db.GetPollIDForOptionID(voteMsg.ID)
		if err != nil {
			server.respondWithError(ERR_BAD_NO_OPTION, err, cVoteTag, writer,
				request)
			return
		}

	default:
		server.respondWithError(ERR_BAD_VOTE_TYPE, err, cVoteTag, writer,
			request)
		return
	}

	// make sure the user is allowed to vote
	if !server.hasPollAccess(user.ID, pollID) {
		server.respondWithError(ERR_ILL_POLL_ACCESS, nil, cVoteTag, writer,
			request)
		return
	}

	// start a transaction
	tx, err := server.db.Begin()
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_TX_BEGIN, err, cVoteTag, writer,
			request)
		return
	}

	// remove all existing votes of the user
	err = database.DeleteVotesForUserTX(user.ID, pollID, tx)
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_DELETE, err, cVoteTag, writer,
			request)
		return
	}

	// if necessary, create a new option
	var optionID int64
	var option polly.Option
	if voteMsg.Type == polly.VOTE_TYPE_UPVOTE {
		optionID = voteMsg.ID
	} else {

		// we have a vote message with type NEW, so we create a new option
		questionID := voteMsg.ID
		option.PollID = pollID
		option.QuestionID = questionID
		option.Value = voteMsg.Value
		err = database.AddOptionTX(&option, tx)
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_ADD, err, cVoteTag, writer,
				request)
			return
		}

		optionID = option.ID
	}

	// insert the vote into the database
	vote := polly.Vote{}
	vote.CreationDate = time.Now().Unix()
	vote.OptionID = optionID
	vote.PollID = pollID
	vote.UserID = user.ID
	err = database.AddVoteTX(&vote, tx)
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_ADD, err, cVoteTag, writer, request)
		return
	}

	// update the poll last updated
	err = database.UpdatePollLastUpdatedTX(pollID, time.Now().Unix(), tx)
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_UPDATE, err, cVoteTag, writer,
			request)
		return
	}

	// update the poll sequence number
	err = database.UpdateSequenceNumberTX(pollID, tx)
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_UPDATE, err, cVoteTag, writer,
			request)
		return
	}

	// retrieve a snapshot of the new poll
	snapshot, err := database.GetPollSnapshotTX(pollID, tx)
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_GET, err, cVoteTag, writer, request)
		return
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		server.respondWithError(ERR_INT_DB_TX_COMMIT, err, cVoteTag, writer,
			request)
		return
	}

	// send a notification to other participants
	err = server.pushClient.NotifyForVote(&server.db, user, optionID,
		voteMsg.Type)
	if err != nil {
		// TODO neaten up
		server.logger.Log(cVoteTag, "Error notifying: "+err.Error(), "::1")
	}

	// construct the response message
	response := polly.VoteResponseMessage{}
	response.Vote = vote
	response.Poll = *snapshot
	if voteMsg.Type == polly.VOTE_TYPE_NEW {
		response.Option = &option
	}

	// send the response message
	responseBody, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cVoteTag, writer,
			request)
		return
	}

	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cVoteTag, writer, request)
		return
	}

}
