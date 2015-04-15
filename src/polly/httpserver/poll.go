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

func (srv *HTTPServer) PostPoll(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(r)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL", fmt.Sprintf("Authentication error: %s",
			err), h)
		w.Header().Set("WWW-authenticate", "Basic")
		http.Error(w, "Authentication error", 401)
		return
	}

	// decode the poll
	var pollMsg PollMessage
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&pollMsg)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL", fmt.Sprintf("Bad JSON: %s", err), h)
		http.Error(w, "Bad JSON.", 400)
		return
	}

	// validate the poll
	if err := isValidPollMessage(srv.db, &pollMsg, usr.Id); err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL", fmt.Sprintf("Invalid poll message: %s",
			err), h)
		http.Error(w, "Bad poll.", 400)
		return
	}

	// insert poll
	pollMsg.MetaData.CreatorId = usr.Id
	pollMsg.Votes = make([]polly.Vote, 0)
	err = srv.InsertPollMessage(&pollMsg)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL", fmt.Sprintf("Database error: %s.", err), h)
		http.Error(w, "Database error.", 500)
		return
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
		http.Error(w, "Marshalling error.", 500)
		return
	}

	// send the response
	_, err = w.Write(responseBody)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL", "Error writing response.", h)
		http.Error(w, "Response error.", 500)
		return
	}

}

func (srv *HTTPServer) GetPoll(w http.ResponseWriter, r *http.Request,
	p httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(r)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Authentication error: %s",
			err), h)
		w.Header().Set("WWW-authenticate", "Basic")
		http.Error(w, "Authentication error", 401)
		return
	}

	// retrieve the poll identifier argument
	pollIdStr := p.ByName(cId)
	pollId, err := strconv.Atoi(pollIdStr)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Bad id: %s", pollIdStr), h)
		http.Error(w, "Bad id.", 400)
		return
	}

	// check whether the user has access rights to the poll
	if !srv.hasPollAccess(usr.Id, pollId) {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL/XX", "User has no access rights to the poll.",
			h)
		http.Error(w, "Illegal operation.", 403)
		return
	}

	// construct the poll message
	pollMsg, err := srv.ConstructPollMessage(pollId)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("No poll with id %s: %s",
			pollIdStr, err), h)
		http.Error(w, "No such poll.", 400)
		return
	}

	// send the response
	responseBody, err := json.MarshalIndent(pollMsg, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/POLL/XX",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
		http.Error(w, "Marshalling error.", 500)
		return
	}

}
