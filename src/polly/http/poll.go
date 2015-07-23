package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	cPostPollTag = "POST/POLL"
	cGetPollTag  = "GET/POLL/XX"
)

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
	if err := isValidPollMessage(server.db, &pollMsg, user.ID); err != nil {
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
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cPostPollTag, err, writer, request)
		return
	}

}

func (server *sServer) GetPoll(writer http.ResponseWriter, request *http.Request,
	params httprouter.Params) {

	// authenticate the user
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cGetPollTag, err, writer, request)
		return
	}

	// retrieve the poll identifier argument
	pollIDStr := params.ByName(cID)
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		server.handleErr(cGetPollTag, cBadIDErr,
			fmt.Sprintf(cLogFmt, cBadIDErr, pollIDStr), 400, writer, request)
		return
	}

	// check whether the user has access rights to the poll
	if !server.hasPollAccess(user.ID, pollID) {
		server.handleIllegalOperation(cGetPollTag, cAccessRightsErr, writer, request)
		return
	}

	// construct the poll message
	pollMsg, err := server.ConstructPollMessage(pollID)
	if err != nil {
		server.handleErr(cGetPollTag, cNoPollErr,
			fmt.Sprintf(cLogFmt, cNoPollErr, pollIDStr), 400, writer, request)
		return
	}

	// send the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleMarshallingError(cGetPollTag, err, writer, request)
		return
	}

}
