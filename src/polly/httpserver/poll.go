package httpserver

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

func (srv *HTTPServer) PostPoll(writer http.ResponseWriter, req *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cPostPollTag, err, writer, req)
		return
	}

	// decode the poll
	var pollMsg PollMessage
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&pollMsg)
	if err != nil {
		srv.handleBadRequest(cPostPollTag, cBadJSONErr, err, writer, req)
		return
	}

	// validate the poll
	if err := isValidPollMessage(srv.db, &pollMsg, usr.Id); err != nil {
		srv.handleBadRequest(cPostPollTag, cBadPollErr, err, writer, req)
		return
	}

	// insert poll
	pollMsg.MetaData.CreatorId = usr.Id
	pollMsg.Votes = make([]polly.Vote, 0)
	err = srv.InsertPollMessage(&pollMsg)
	if err != nil {
		srv.handleDatabaseError(cPostPollTag, err, writer, req)
		return
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	if err != nil {
		srv.handleMarshallingError(cPostPollTag, err, writer, req)
		return
	}

	// send the response
	_, err = writer.Write(responseBody)
	if err != nil {
		srv.handleWritingError(cPostPollTag, err, writer, req)
		return
	}

}

func (srv *HTTPServer) GetPoll(writer http.ResponseWriter, req *http.Request,
	params httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cGetPollTag, err, writer, req)
		return
	}

	// retrieve the poll identifier argument
	pollIdStr := params.ByName(cId)
	pollId, err := strconv.Atoi(pollIdStr)
	if err != nil {
		srv.handleErr(cGetPollTag, cBadIdErr,
			fmt.Sprintf(cLogFmt, cBadIdErr, pollIdStr), 400, writer, req)
		return
	}

	// check whether the user has access rights to the poll
	if !srv.hasPollAccess(usr.Id, pollId) {
		srv.handleIllegalOperation(cGetPollTag, cAccessRightsErr, writer, req)
		return
	}

	// construct the poll message
	pollMsg, err := srv.ConstructPollMessage(pollId)
	if err != nil {
		srv.handleErr(cGetPollTag, cNoPollErr,
			fmt.Sprintf(cLogFmt, cNoPollErr, pollIdStr), 400, writer, req)
		return
	}

	// send the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		srv.handleMarshallingError(cGetPollTag, err, writer, req)
		return
	}

}
