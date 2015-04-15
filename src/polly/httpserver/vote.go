package httpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"polly"
	"polly/database"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	VOTE_TYPE_NEW    = 0
	VOTE_TYPE_UPVOTE = 1
)

type VoteMessage struct {
	Type  int    `json:"type"`
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type VoteResponseMessage struct {
	Option *polly.Option `json:"option,omitempty"`
	Vote   *polly.Vote   `json:"vote"`
}

func (srv *HTTPServer) Vote(w http.ResponseWriter, r *http.Request,
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

	// decode the vote message
	var voteMsg VoteMessage
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&voteMsg)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/POLL", fmt.Sprintf("Bad JSON: %s", err), h)
		http.Error(w, "Bad JSON.", 400)
		return
	}

	// retrieve the poll id belonging to the option or question id
	var pollId int
	switch voteMsg.Type {
	case VOTE_TYPE_NEW:
		pollId, err = srv.db.PollIdForQuestionId(voteMsg.Id)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/VOTE", fmt.Sprintf(
				"Unknown question id: %d: %s.", voteMsg.Id, err), h)
			http.Error(w, "Unknown question id.", 400)
			return
		} else if len(voteMsg.Value) == 0 {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/VOTE",
				"Invalid vote message: empty value field for vote message "+
					"with type NEW.", h)
			http.Error(w, "Bad vote message.", 400)
			return
		}
	case VOTE_TYPE_UPVOTE:
		pollId, err = srv.db.PollIdForOptionId(voteMsg.Id)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/VOTE", fmt.Sprintf("Unknown option id: %d.",
				voteMsg.Id), h)
			http.Error(w, fmt.Sprintf("Unknown option id: %d: %s.", voteMsg.Id,
				err), 400)
			return
		}
	default:
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Bad vote type: %d.",
			voteMsg.Type), h)
		http.Error(w, "Bad vote type.", 400)
		return
	}

	// make sure the user is allowed to vote
	if !srv.hasPollAccess(usr.Id, pollId) {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", "User has no access rights to the poll.", h)
		http.Error(w, "Illegal operation.", 403)
		return
	}

	// start a transaction
	transaction, err := srv.db.Begin()
	if err != nil {
		transaction.Rollback()
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err), h)
		http.Error(w, "Database error", 500)
		return
	}

	// remove all existing votes of the user
	err = database.DelVotesForUserTx(usr.Id, pollId, transaction)
	if err != nil {
		transaction.Rollback()
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err), h)
		http.Error(w, "Database error", 500)
		return
	}

	// if necessary, create a new option
	var optionId int
	var option polly.Option
	if voteMsg.Type == VOTE_TYPE_UPVOTE {
		optionId = voteMsg.Id
	} else {

		// we have a vote message with type NEW, so we create a new option
		questionId := voteMsg.Id
		option.PollId = pollId
		option.QuestionId = questionId
		option.Value = voteMsg.Value
		err = srv.db.AddOptionTx(&option, transaction)
		if err != nil {
			transaction.Rollback()
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err),
				h)
			http.Error(w, "Database error", 500)
			return
		}

		optionId = option.Id
	}

	// insert the vote into the database
	vote := polly.Vote{}
	vote.CreationDate = time.Now().Unix()
	vote.OptionId = optionId
	vote.PollId = pollId
	vote.UserId = usr.Id
	err = database.AddVoteTx(&vote, transaction)
	if err != nil {
		transaction.Rollback()
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err), h)
		http.Error(w, "Database error", 500)
		return
	}

	// update the poll last updated
	err = database.UpdatePollLastUpdatedTx(pollId, time.Now().Unix(),
		transaction)
	if err != nil {
		transaction.Rollback()
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err), h)
		http.Error(w, "Database error", 500)
		return
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err), h)
		http.Error(w, "Database error", 500)
		return
	}

	// construct the response message
	response := VoteResponseMessage{}
	response.Vote = &vote
	if voteMsg.Type == VOTE_TYPE_NEW {
		response.Option = &option
	}

	// send the response message
	responseBody, err := json.MarshalIndent(response, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/VOTE",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
		http.Error(w, "Marshalling error.", 500)
		return
	}

}
