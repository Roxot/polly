package http

import (
	"encoding/json"
	"net/http"
	"polly"

	"polly/internal/github.com/julienschmidt/httprouter"
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
