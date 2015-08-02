package http

import (
	"encoding/json"
	"net/http"
	"polly"
	"strconv"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cPostPollTag    = "POST/POLL"
	cGetPollBulkTag = "GET/POLLS"
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
