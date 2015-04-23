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
	cListUserPollsTag = "USER/POLLS"
)

type PollSnapshot struct {
	PollId      int `db:"poll_id" json:"poll_id"`
	LastUpdated int `db:"last_updated" json:"last_updated"`
}

type PollList struct {
	Snapshots  []polly.PollSnapshot `json:"polls"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	NumResults int                  `json:"num_results"`
	Total      int64                `json:"total"`
}

func (srv *HTTPServer) ListUserPolls(writer http.ResponseWriter, req *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cListUserPollsTag, err, writer, req)
		return
	}

	// retrieve the page argument
	var page int
	pageStrings := req.URL.Query()[cPage]
	if len(pageStrings) > 0 {

		// convert the page argument to an integer
		pageStr := pageStrings[0]
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			srv.handleErr(cListUserPollsTag, cBadPageErr,
				fmt.Sprintf(cLogFmt, cBadPageErr, pageStr), 400, writer, req)
			return
		}

	} else {
		page = 1
	}

	// retrieve poll snapshots
	offset := (page - 1) * cPollListMax
	snapshots, err := srv.db.PollSnapshotsByUserId(usr.Id, cPollListMax,
		offset)
	if err != nil {
		srv.handleDatabaseError(cListUserPollsTag, err, writer, req)
		return
	}

	// construct the PollList object
	pollList := PollList{}
	pollList.Snapshots = snapshots
	pollList.Page = page
	pollList.PageSize = cPollListMax
	pollList.NumResults = len(snapshots)
	pollList.Total = srv.db.CountPollsForUser(usr.Id)

	// send the response
	responseBody, err := json.MarshalIndent(pollList, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		srv.handleMarshallingError(cListUserPollsTag, err, writer, req)
		return
	}
}
