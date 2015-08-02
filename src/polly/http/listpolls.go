package http

import (
	"encoding/json"
	"net/http"
	"polly"
	"strconv"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cListUserPollsTag = "GET/LIST_POLLS"
)

func (server *sServer) ListPolls(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {
	var err error

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cListUserPollsTag, writer,
			request)
		return
	}

	// retrieve the page argument
	var page int
	pageStrings := request.URL.Query()[cPage]
	if len(pageStrings) > 0 {

		// convert the page argument to an integer
		pageStr := pageStrings[0]
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			server.respondWithError(ERR_BAD_PAGE, err, cListUserPollsTag,
				writer, request)
			return
		}

	} else {
		page = 1
	}

	// retrieve poll snapshots
	offset := (page - 1) * cPollListMax
	snapshots, err := server.db.GetPollSnapshotsByUserID(user.ID, cPollListMax,
		offset)
	if err != nil {
		server.respondWithError(ERR_INT_DB_GET, err, cListUserPollsTag, writer,
			request)
		return
	}

	// construct the PollList object
	pollListMsg := polly.PollListMessage{}
	pollListMsg.Snapshots = snapshots
	pollListMsg.Page = page
	pollListMsg.PageSize = cPollListMax
	pollListMsg.NumResults = len(snapshots)
	pollListMsg.Total = server.db.CountPollsForUser(user.ID)

	// marshall the response
	responseBody, err := json.MarshalIndent(pollListMsg, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cListUserPollsTag,
			writer, request)
		return
	}

	// send the response
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cListUserPollsTag, writer,
			request)
		return
	}
}
