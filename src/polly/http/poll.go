package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly"
	"strconv"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cPostPollTag    = "POST/POLL"
	cGetPollBulkTag = "GET/POLLS"
)

type PollBulk struct { // TODO private or model
	Polls []PollMessage `json:"polls"`
}

func (server *sServer) PostPoll(writer http.ResponseWriter, request *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cPostPollTag, err, writer, request)
		return
	}

	// decode the poll
	var pollMsg PollMessage
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&pollMsg)
	if err != nil {
		server.handleBadRequest(cPostPollTag, cBadJSONErr, err, writer, request)
		return
	}

	// validate the poll
	if err := isValidPollMessage(&server.db, &pollMsg, user.ID); err != nil {
		server.handleBadRequest(cPostPollTag, cBadPollErr, err, writer, request)
		return
	}

	// insert poll
	pollMsg.MetaData.CreatorID = user.ID
	pollMsg.Votes = make([]polly.Vote, 0)
	err = server.InsertPollMessage(&pollMsg)
	if err != nil {
		server.handleDatabaseError(cPostPollTag, err, writer, request)
		return
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	if err != nil {
		server.handleMarshallingError(cPostPollTag, err, writer, request)
		return
	}

	// send the response
	SetJSONContentType(writer)
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cPostPollTag, err, writer, request)
		return
	}

}

func (server *sServer) GetPollBulk(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	// authenticate the request
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cGetPollBulkTag, err, writer, request)
		return
	}

	// retrieve the list of identifiers
	ids := request.URL.Query()[cID]
	if len(ids) > cBulkPollMax {
		server.handleErr(cGetPollBulkTag, cIDListLengthErr,
			fmt.Sprintf("%s: %d", cIDListLengthErr, len(ids)), 400, writer,
			request)
		return
	}

	// construct the PollBulk object
	pollBulk := PollBulk{}
	pollBulk.Polls = make([]PollMessage, len(ids))
	for idx, idString := range ids {

		// convert the id to an integer
		id, err := strconv.Atoi(idString)
		if err != nil {
			server.handleBadRequest(cGetPollBulkTag, cBadIDErr, err, writer,
				request)
			return
		}

		// make sure the user is authorized to receive the poll
		if !server.hasPollAccess(user.ID, id) {
			server.handleIllegalOperation(cGetPollBulkTag, cAccessRightsErr,
				writer, request)
			return
		}

		// construct the poll message
		pollMsg, err := server.ConstructPollMessage(id)
		if err != nil {
			server.handleErr(cGetPollBulkTag, cNoPollErr, cNoPollErr, 400,
				writer, request)
			return
		}

		pollBulk.Polls[idx] = *pollMsg
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(pollBulk, "", "\t")
	if err != nil {
		server.handleMarshallingError(cGetPollBulkTag, err, writer, request)
		return
	}

	// send a 200 OK response
	SetJSONContentType(writer)
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cGetPollBulkTag, err, writer, request)
		return
	}
}
