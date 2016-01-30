package http

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"github.com/roxot/polly"
	"github.com/roxot/polly/database"
)

const (
	cGetUserBulkTag = "GET/USERS"
	cUpdateUserTag  = "PUT/USER"
	cAddUserTag     = "POST/ADDUSER"
)

func (server *sServer) GetUserBulk(writer http.ResponseWriter,
	request *http.Request, params httprouter.Params) {
	var err error

	// authenticate the user
	_, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cGetUserBulkTag, writer, request)
		return
	}

	// retrieve the list of identifiers
	ids := request.URL.Query()[cID]
	if len(ids) > cBulkUserMax {
		server.respondWithError(ERR_ILL_TOO_MANY_IDS, nil, cGetUserBulkTag,
			writer, request)
		return
	}

	// construct the UserBulk object
	userBulkMsg := polly.UserBulkMessage{}
	userBulkMsg.Users = make([]polly.PublicUser, len(ids))
	idx := -1
	for _, idString := range ids {

		// convert the id to an integer
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			server.respondWithError(ERR_BAD_ID, err, cGetUserBulkTag, writer,
				request)
			return
		}

		// retrieve the user object
		user, err := server.db.GetPublicUserByID(id)
		if err == nil {
			idx += 1
			userBulkMsg.Users[idx] = *user
		}

	}

	// Only include users that were actually found
	userBulkMsg.Users = userBulkMsg.Users[:idx+1]

	// marshall the response
	responseBody, err := json.MarshalIndent(userBulkMsg, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cGetUserBulkTag, writer,
			request)
		return
	}

	// send the response
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cGetUserBulkTag, writer,
			request)
		return
	}

}

func (server *sServer) UpdateUser(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {
	var err error

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cUpdateUserTag, writer, request)
		return
	}

	// decode the given user
	var updateUserMsg polly.UpdateUserMessage
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&updateUserMsg)
	if err != nil {
		server.respondWithError(ERR_BAD_JSON, err, cUpdateUserTag, writer,
			request)
		return
	}

	// update display name
	if updateUserMsg.DeviceGUID != nil {
		user.DeviceGUID = *(updateUserMsg.DeviceGUID)
		err = server.db.UpdateDeviceGUID(user.ID, user.DeviceGUID)
		if err != nil {
			server.respondWithError(ERR_INT_DB_UPDATE, err, cUpdateUserTag,
				writer, request)
			return
		}
	}

	// update device GUID
	if updateUserMsg.DisplayName != nil {
		user.DisplayName = *(updateUserMsg.DisplayName)
		err = server.db.UpdateDisplayName(user.ID, user.DisplayName)
		if err != nil {
			server.respondWithError(ERR_INT_DB_UPDATE, err, cUpdateUserTag,
				writer, request)
			return
		}
	}

	// create the response body
	responseBody, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cUpdateUserTag,
			writer, request)
		return
	}

	// send the user a 200 OK with updated user
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cUpdateUserTag, writer,
			request)
		return
	}

}

// POST /api/v0.1/adduser.json
func (server *sServer) AddUser(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cAddUserTag, writer, request)
		return
	}

	// decode the given user
	var addUserMsg polly.AddUserMessage
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&addUserMsg)
	if err != nil {
		server.respondWithError(ERR_BAD_JSON, err, cAddUserTag, writer, request)
		return
	}

	// retrieve the creator of the poll
	creatorID, err := server.db.GetPollCreatorID(addUserMsg.PollID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cAddUserTag, writer,
			request)
		return
	}

	// make sure the user is the creator of the poll
	if creatorID != user.ID {
		server.respondWithError(ERR_ILL_NOT_CREATOR, nil, cAddUserTag, writer,
			request)
		return
	}

	// retrieve the closing date
	closingDate, err := server.db.GetClosingDate(addUserMsg.PollID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cAddUserTag, writer,
			request)
		return
	}

	// make sure the poll hasn't closed yet
	currentTime := time.Now().UnixNano() / 1000000
	if currentTime > closingDate {
		server.respondWithError(ERR_ILL_POLL_CLOSED, nil, cAddUserTag, writer,
			request)
		return
	}

	// check if the new user exists
	newUser, err := server.db.GetUserByID(addUserMsg.User.ID)
	if err != nil {
		server.respondWithError(ERR_BAD_NO_USER, err, cAddUserTag, writer,
			request)
		return
	}

	// get the corresponding question
	question, err := server.db.GetQuestionByPollID(addUserMsg.PollID)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cAddUserTag, writer,
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
			server.respondWithError(ERR_INT_DB_TX_BEGIN, err, cAddUserTag,
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
		err = database.UpdatePollTX(addUserMsg.PollID, currentTime,
			polly.EVENT_TYPE_NEW_PARTICIPANT, newUser.DisplayName, question.Title, tx)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cAddUserTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {

				tx.Rollback()
				server.respondWithError(ERR_INT_DB_UPDATE, err, cAddUserTag,
					writer, request)
				return
			}
		}

		// make sure the user is not already in the poll
		isParticipant, err := database.ExistsParticipantTX(addUserMsg.User.ID,
			addUserMsg.PollID, tx)
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_GET, err, cAddUserTag, writer,
				request)
			return
		} else if isParticipant {
			tx.Rollback()
			server.respondWithError(ERR_BAD_DUPLICATE_PARTICIPANT, nil,
				cAddUserTag, writer, request)
			return
		}

		// add the user to the poll
		newParticipant := &polly.Participant{PollID: addUserMsg.PollID,
			UserID: addUserMsg.User.ID}
		err = database.AddParticipantTX(newParticipant, tx)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok &&
				pqErr.Code == database.ERR_SERIALIZATION_FAILURE {
				server.logger.Log(cAddUserTag, fmt.Sprintf("%d: %s",
					transactionNumber, "Serialization failure, retrying..."),
					"::1")
				continue
			} else {
				tx.Rollback()
				server.respondWithError(ERR_INT_DB_ADD, err, cAddUserTag, writer,
					request)
				return
			}
		}

		// commit the transaction
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			server.respondWithError(ERR_INT_DB_TX_COMMIT, err, cAddUserTag,
				writer, request)
			return
		}

		retryTransaction = false
	}

	// notify the users of the poll of the change
	err = server.pushClient.NotifyForNewParticipant(&server.db, user,
		addUserMsg.PollID, question.Title, newUser)
	if err != nil {
		// TODO neaten up
		server.logger.Log(cAddUserTag, "Error notifying: "+err.Error(), "::1")
	}

	// respond with 200 OK
	server.respondOkay(writer, request)
}
