package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	cGetPollBulkTag = "GET/POLL"
)

type PollBulk struct {
	Polls []PollMessage `json:"polls"`
}

func (srv *HTTPServer) GetPollBulk(writer http.ResponseWriter,
	req *http.Request, _ httprouter.Params) {

	// authenticate the request
	usr, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cGetPollBulkTag, err, writer, req)
		return
	}

	// retrieve the list of identifiers
	ids := req.URL.Query()[cId]
	if len(ids) > cBulkPollMax {
		srv.handleErr(cGetPollBulkTag, cIdListLengthErr,
			fmt.Sprintf("%s: %d", cIdListLengthErr, len(ids)), 400, writer, req)
		return
	}

	// construct the PollBulk object
	pollBulk := PollBulk{}
	pollBulk.Polls = make([]PollMessage, len(ids))
	for idx, idString := range ids {

		// convert the id to an integer
		id, err := strconv.Atoi(idString)
		if err != nil {
			srv.handleBadRequest(cGetPollBulkTag, cBadIdErr, err, writer, req)
			return
		}

		// make sure the user is authorized to receive the poll
		if !srv.hasPollAccess(usr.Id, id) {
			srv.handleIllegalOperation(cGetPollBulkTag, cAccessRightsErr,
				writer, req)
			return
		}

		// construct the poll message
		pollMsg, err := srv.ConstructPollMessage(id)
		if err != nil {
			srv.handleErr(cGetPollBulkTag, cNoPollErr, cNoPollErr, 400,
				writer, req)
			return
		}

		pollBulk.Polls[idx] = *pollMsg
	}

	responseBody, err := json.MarshalIndent(pollBulk, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		srv.handleWritingError(cGetPollBulkTag, err, writer, req)
		return
	}
}
