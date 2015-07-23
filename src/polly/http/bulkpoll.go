package http

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

	responseBody, err := json.MarshalIndent(pollBulk, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cGetPollBulkTag, err, writer, request)
		return
	}
}
