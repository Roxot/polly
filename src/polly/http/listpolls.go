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
	cListUserPollsTag = "GET/LIST_POLLS"
)

type PollList struct {
	Snapshots  []polly.PollSnapshot `json:"polls"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	NumResults int                  `json:"num_results"`
	Total      int64                `json:"total"`
}

func (server *sServer) ListPolls(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	// authenticate the user
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cListUserPollsTag, err, writer, request)
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
			server.handleErr(cListUserPollsTag, cBadPageErr,
				fmt.Sprintf(cLogFmt, cBadPageErr, pageStr), 400, writer,
				request)
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
		server.handleDatabaseError(cListUserPollsTag, err, writer, request)
		return
	}

	// construct the PollList object
	pollList := PollList{}
	pollList.Snapshots = snapshots
	pollList.Page = page
	pollList.PageSize = cPollListMax
	pollList.NumResults = len(snapshots)
	pollList.Total = server.db.CountPollsForUser(user.ID)

	// marshall the response
	responseBody, err := json.MarshalIndent(pollList, "", "\t")
	if err != nil {
		server.handleMarshallingError(cListUserPollsTag, err, writer, request)
		return
	}

	// send the response
	SetJSONContentType(writer)
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleMarshallingError(cListUserPollsTag, err, writer, request)
		return
	}
}
