package http

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/roxot/polly"
	"github.com/roxot/polly/database"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
)

const (
	cVoteTag     = "POST/VOTE"
	cUndoVoteTag = "DELETE/VOTE"
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
	var optionTitle string
	switch voteMsg.Type {
	case polly.VOTE_TYPE_NEW:
		question, err := server.db.GetQuestionByID(voteMsg.ID)
		if err != nil {
			server.respondWithError(ERR_BAD_NO_QUESTION, err, cVoteTag, writer,
				request)
			return

		}

		pollID = question.ID
		if question.Type != polly.QUESTION_TYPE_OPEN &&
			question.Type != polly.QUESTION_TYPE_MOVIE_OPEN {

			server.respondWithError(ERR_ILL_ADD_OPTION, nil, cVoteTag, writer,
				request)
			return
		} else if len(voteMsg.Value) == 0 {
			server.respondWithError(ERR_BAD_EMPTY_OPTION, nil, cVoteTag, writer,
				request)
			return
		}

		optionTitle = voteMsg.Value
	case polly.VOTE_TYPE_UPVOTE:
		option, err := server.db.GetOptionByID(voteMsg.ID)
		if err != nil {
			server.respondWithError(ERR_BAD_NO_OPTION, err, cVoteTag, writer,
				request)
			return
		}

		pollID, optionTitle = option.PollID, option.Value
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

	// retrieve the closing date
	closingDate, err := server.db.GetClosingDate(pollID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cVoteTag, writer,
			request)
		return
	}

	// make sure the poll hasn't closed yet
	currentTime := time.Now().UnixNano() / 1000000
	if currentTime > closingDate {
		server.respondWithError(ERR_ILL_POLL_CLOSED, nil, cVoteTag, writer,
			request)
		return
	}

	var optionID int64
	var option polly.Option
	var snapshot *polly.PollSnapshot
	var vote polly.Vote
	retryTransaction := true
	transactionNumber := rand.Int()
	for retryTransaction {

		// start a transaction
		tx, err := server.db.Begin()
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_BEGIN, err, cVoteTag, writer,
				request)
			return
		}

		// set the transaction isolation level
		_, err = tx.Exec("set transaction isolation level serializable;")
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_SET_TX_LEVEL, err,
				cVoteTag, writer, request)
			return
		}

		// update the poll last updated and seq number
		// user, optionID, voteMsg.Type
		err = database.UpdatePollTX(pollID, currentTime, voteMsg.Type,
			user.DisplayName, user.ID, optionTitle, tx)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cVoteTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {

				tx.Rollback()
				server.respondWithError(ERR_INT_DB_UPDATE, err, cVoteTag,
					writer, request)
				return
			}
		}

		// retrieve a snapshot of the new poll
		snapshot, err = database.GetPollSnapshotTX(pollID, tx)
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_GET, err, cVoteTag, writer, request)
			return
		}

		// remove all existing votes of the user
		err = database.DeleteVotesForUserTX(user.ID, pollID, tx)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cVoteTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {
				tx.Rollback()
				server.respondWithError(ERR_INT_DB_DELETE, err, cVoteTag, writer,
					request)
				return
			}
		}

		// if necessary, create a new option, otherwise update the existing
		// option its sequence number
		if voteMsg.Type == polly.VOTE_TYPE_UPVOTE {
			optionID = voteMsg.ID
			err := database.UpdateOptionSequenceNumberTX(optionID,
				snapshot.SequenceNumber, tx)
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok &&
					pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
					server.logger.Log(cVoteTag, fmt.Sprintf("%d: %s",
						transactionNumber,
						"Serialization failure, retrying..."), "::1")
					continue
				} else {
					tx.Rollback()
					server.respondWithError(ERR_INT_DB_UPDATE, err, cVoteTag,
						writer, request)
					return
				}
			}

		} else {

			// we have a vote message with type NEW, so we create a new option
			questionID := voteMsg.ID
			option.PollID = pollID
			option.QuestionID = questionID
			option.Value = voteMsg.Value
			option.SequenceNumber = snapshot.SequenceNumber
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
		vote = polly.Vote{}
		vote.CreationDate = currentTime
		vote.OptionID = optionID
		vote.PollID = pollID
		vote.UserID = user.ID
		err = database.AddVoteTX(&vote, tx)
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_ADD, err, cVoteTag, writer, request)
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

		retryTransaction = false
	}

	// send a notification to other participants
	err = server.pushClient.NotifyForVote(&server.db, user, optionTitle, pollID,
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

	// marshal the response body
	responseBody, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cVoteTag, writer,
			request)
		return
	}

	// send the response message
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cVoteTag, writer, request)
		return
	}

}

func (server *sServer) UndoVote(writer http.ResponseWriter,
	request *http.Request, params httprouter.Params) {

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cUndoVoteTag, writer, request)
		return
	}

	// convert the id to an integer and make sure its a valid integer value
	ids := request.URL.Query()[cID]
	if len(ids) == 0 {
		server.respondWithError(ERR_BAD_NO_ID, nil, cUndoVoteTag, writer,
			request)
		return
	}

	// parse the provided id to an integer
	id, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		server.respondWithError(ERR_BAD_ID, err, cUndoVoteTag, writer, request)
		return
	}

	// retrieve the vote object
	vote, err := server.db.GetVoteByID(id)
	if err != nil {
		server.respondWithError(ERR_BAD_NO_VOTE, err, cUndoVoteTag, writer,
			request)
		return
	}

	// retrieve the closing date
	closingDate, err := server.db.GetClosingDate(vote.PollID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cUndoVoteTag, writer,
			request)
		return
	}

	// make sure the poll hasn't closed yet
	currentTime := time.Now().UnixNano() / 1000000
	if currentTime > closingDate {
		server.respondWithError(ERR_ILL_POLL_CLOSED, nil, cUndoVoteTag, writer,
			request)
		return
	}

	// retrieve the belonging option object
	option, err := server.db.GetOptionByID(vote.OptionID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cUndoVoteTag, writer, request)
		return
	}

	var snapshot *polly.PollSnapshot
	retryTransaction := true
	transactionNumber := rand.Int()
	for retryTransaction {

		// start a transaction
		tx, err := server.db.Begin()
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_BEGIN, err, cUndoVoteTag,
				writer, request)
			return
		}

		// set the transaction isolation level
		_, err = tx.Exec("set transaction isolation level serializable;")
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_SET_TX_LEVEL, err,
				cUndoVoteTag, writer, request)
			return
		}

		// update the poll last updated and seq number
		err = database.UpdatePollTX(vote.PollID, currentTime,
			polly.EVENT_TYPE_UNDONE_VOTE, user.DisplayName, user.ID, option.Value, tx)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cUndoVoteTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {

				tx.Rollback()
				server.respondWithError(ERR_INT_DB_UPDATE, err, cUndoVoteTag,
					writer, request)
				return
			}
		}

		// retrieve a snapshot of the new poll
		snapshot, err = database.GetPollSnapshotTX(vote.PollID, tx)
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_GET, err, cUndoVoteTag, writer,
				request)
			return
		}

		// delete the vote
		err = server.db.DeleteVoteByIDForUser(id, user.ID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cUndoVoteTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {
				tx.Rollback()
				server.respondWithError(ERR_BAD_NO_VOTE, err, cUndoVoteTag,
					writer, request)
				return
			}
		}

		// commit the transaction
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_COMMIT, err, cUndoVoteTag,
				writer, request)
			return
		}

		retryTransaction = false
	}

	// notify the poll participants
	err = server.pushClient.NotifyForUndoneVote(&server.db, user, option.Value,
		vote.PollID)
	if err != nil {
		// TODO neaten up
		server.logger.Log(cUndoVoteTag, "Error notifying: "+err.Error(), "::1")
	}

	// marshal the response body
	responseBody, err := json.MarshalIndent(snapshot, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cUndoVoteTag, writer,
			request)
		return
	}

	// send the response message
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cUndoVoteTag, writer,
			request)
		return
	}
}
