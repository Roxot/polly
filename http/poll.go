package http

import (
	"encoding/json"
	"fmt"
	"github.com/roxot/polly"
	"github.com/roxot/polly/database"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
)

const (
	cPostPollTag    = "POST/POLL"
	cGetPollBulkTag = "GET/POLLS"
	cLeavePollTag   = "DELETE/POLL"
)

func (server *sServer) PostPoll(writer http.ResponseWriter, request *http.Request,
	_ httprouter.Params) {
	var err error

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cPostPollTag, writer, request)
		return
	}

	// decode the poll
	var pollMsg polly.PollMessage
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&pollMsg)
	if err != nil {
		server.respondWithError(ERR_BAD_JSON, err, cPostPollTag, writer,
			request)
		return
	}

	// validate the poll
	if errCode = isValidPollMessage(&server.db, &pollMsg, user.ID); errCode !=
		NO_ERR {
		server.respondWithError(errCode, nil, cPostPollTag, writer, request)
		return
	}

	// insert poll
	pollMsg.MetaData.CreatorID = user.ID
	pollMsg.Votes = make([]polly.Vote, 0)
	err = server.db.InsertPollMessage(&pollMsg)
	if err != nil {
		server.respondWithError(ERR_INT_DB_ADD, err, cPostPollTag, writer,
			request)
		return
	}

	// notify the poll participants of the creation of the poll
	err = server.pushClient.NotifyForNewPoll(&server.db, user,
		pollMsg.MetaData.ID, pollMsg.Question.Title)
	if err != nil {
		// TODO neaten up
		server.logger.Log(cPostPollTag, "Error notifying: "+err.Error(), "::1")
	}

	// schedule the closing of the poll
	closingDate := time.Unix(0, 1000000*pollMsg.MetaData.ClosingDate)
	pollToClose := tPollToClose{pollMsg.MetaData.ID, pollMsg.Question.Title}
	_, err = server.cpScheduler.Schedule(0, closingDate, &pollToClose)
	if err != nil {
		server.respondWithError(ERR_INT_CP_SCHEDULER, err, cPostPollTag, writer,
			request)
		return
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cPostPollTag, writer,
			request)
		return
	}

	// send the response
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cPostPollTag, writer,
			request)
	}
}

func (server *sServer) GetPollBulk(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	// authenticate the request
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cGetPollBulkTag, writer, request)
		return
	}

	// retrieve the list of identifiers
	ids := request.URL.Query()[cID]
	if len(ids) > cBulkPollMax {
		server.respondWithError(ERR_ILL_TOO_MANY_IDS, nil, cGetPollBulkTag,
			writer, request)
		return
	}

	// construct the PollBulk object
	pollBulkMsg := polly.PollBulkMessage{}
	pollBulkMsg.Polls = make([]polly.PollMessage, len(ids))
	for idx, idString := range ids {

		// convert the id to an integer
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			server.respondWithError(ERR_BAD_ID, err, cGetPollBulkTag, writer,
				request)
			return
		}

		// make sure the user is authorized to receive the poll
		if !server.hasPollAccess(user.ID, id) {
			server.respondWithError(ERR_ILL_POLL_ACCESS, nil, cGetPollBulkTag,
				writer, request)
			return
		}

		// construct the poll message
		pollMsg, err := server.db.ConstructPollMessage(id)
		if err != nil {
			server.respondWithError(ERR_BAD_NO_POLL, err, cGetPollBulkTag,
				writer, request)
			return
		}

		pollBulkMsg.Polls[idx] = *pollMsg
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(pollBulkMsg, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cGetPollBulkTag, writer,
			request)
		return
	}

	// send a 200 OK response
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cGetPollBulkTag, writer,
			request)
		return
	}
}

func (server *sServer) LeavePoll(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	// authenticate the request
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cLeavePollTag, writer, request)
		return
	}

	// convert the id to an integer and make sure it's a valid integer value
	ids := request.URL.Query()[cID]
	if len(ids) == 0 {
		server.respondWithError(ERR_BAD_NO_ID, nil, cLeavePollTag, writer,
			request)
		return
	}

	// parse the provided poll id to an integer
	pollID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		server.respondWithError(ERR_BAD_ID, err, cLeavePollTag, writer, request)
		return
	}

	// retrieve the closing date
	closingDate, err := server.db.GetClosingDate(pollID)
	if err != nil {
		server.respondWithError(ERR_BAD_NO_POLL, err, cLeavePollTag, writer,
			request)
		return
	}

	// make sure the poll hasn't closed yet
	currentTime := time.Now().UnixNano() / 1000000
	if currentTime > closingDate {
		server.respondWithError(ERR_ILL_POLL_CLOSED, nil, cLeavePollTag, writer,
			request)
		return
	}

	// retrieve the poll question
	question, err := server.db.GetQuestionByPollID(pollID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cLeavePollTag, writer,
			request)
		return
	}

	retryTransaction := true
	transactionNumber := rand.Int()
	for retryTransaction {

		// start a transaction
		tx, err := server.db.Begin()
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_BEGIN, err, cLeavePollTag,
				writer, request)
			return
		}

		// set the transaction isolation level
		_, err = tx.Exec("set transaction isolation level serializable;")
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		// update the poll last updated and seq number
		err = database.UpdatePollTX(pollID, currentTime, tx)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cLeavePollTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {

				tx.Rollback()
				server.respondWithError(ERR_INT_DB_UPDATE, err, cLeavePollTag,
					writer, request)
				return
			}
		}

		// delete the participant from the poll
		err = server.db.DeleteParticipant(user.ID, pollID)
		if err != nil { // TODO what if internal?
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cLeavePollTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {
				tx.Rollback()
				server.respondWithError(ERR_BAD_NO_POLL, err, cLeavePollTag, writer,
					request)
				return
			}
		}

		// commit the transaction
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_COMMIT, err, cLeavePollTag,
				writer, request)
			return
		}

		retryTransaction = false
	}

	// notify the poll participants
	err = server.pushClient.NotifyForParticipantLeft(&server.db, user, pollID,
		question.Title)
	if err != nil {
		// TODO neaten up
		server.logger.Log(cLeavePollTag, "Error notifying: "+err.Error(), "::1")
	}

	// respond with 200 ok
	server.respondOkay(writer, request)
}
