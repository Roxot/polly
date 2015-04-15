package httpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PollBulk struct {
	Polls []PollMessage `json:"polls"`
}

func (srv *HTTPServer) GetPollBulk(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	// authenticate the request
	usr, err := srv.authenticateRequest(r)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL", fmt.Sprintf("Authentication error: %s",
			err), h)
		w.Header().Set("WWW-authenticate", "Basic")
		http.Error(w, "Authentication error", 401)
		return
	}

	// retrieve the list of identifiers
	ids := r.URL.Query()[cId]
	if len(ids) > cBulkPollMax {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL", fmt.Sprintf(
			"Id list longer than limit: %d > %d", len(ids), cBulkPollMax), h)
		http.Error(w, "Id list longer than limit", 400)
		return
	}

	// construct the PollBulk object
	pollBulk := PollBulk{}
	pollBulk.Polls = make([]PollMessage, len(ids))
	for idx, idString := range ids {

		// convert the id to an integer
		id, err := strconv.Atoi(idString)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("GET/POLL", fmt.Sprintf("Bad id: %s", idString), h)
			http.Error(w, "Bad id.", 400)
			return
		}

		// make sure the user is authorized to receive the poll
		if !srv.hasPollAccess(usr.Id, id) {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("GET/POLL", "User has no access rights to the poll.",
				h)
			http.Error(w, "Illegal operation.", 403)
			return
		}

		// construct the poll message
		pollMsg, err := srv.ConstructPollMessage(id)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("GET/POLL", fmt.Sprintf("No poll with id %s: %s",
				idString, err), h)
			http.Error(w, "No such poll.", 400)
			return
		}

		pollBulk.Polls[idx] = *pollMsg
	}

	responseBody, err := json.MarshalIndent(pollBulk, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
		http.Error(w, "Marshalling error.", 500)
	}
}
