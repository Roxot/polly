package httpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"polly"
	"strconv"

	"github.com/julienschmidt/httprouter"
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

func (srv *HTTPServer) ListUserPolls(w http.ResponseWriter, r *http.Request,
	p httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(r)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("USER/POLLS", fmt.Sprintf("Authentication error: %s",
			err), h)
		w.Header().Set("WWW-authenticate", "Basic")
		http.Error(w, "Authentication error", 401)
		return
	}

	// retrieve the page argument
	var page int
	pageStrings := r.URL.Query()[cPage]
	if len(pageStrings) > 0 {

		// convert the page argument to an integer
		pageStr := pageStrings[0]
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("USER/POLLS", fmt.Sprintf("Bad page: %s", pageStr),
				h)
			http.Error(w, "Bad page.", 400)
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
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("USER/POLLS", fmt.Sprintf("No polls for user: %s", err),
			h)
		http.Error(w, "Database error.", 500) // TODO when does this happen?
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
	_, err = w.Write(responseBody)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("USER/POLLS",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
		http.Error(w, "Marshalling error.", 500)
	}
}
