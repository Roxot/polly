package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly"
	"polly/database"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (srv *HTTPServer) Vote(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(r)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Authentication error: %s",
			err))
		http.Error(w, "Authentication error", 401)
		return
	}

	// decode the vote
	var vote polly.Vote
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&vote)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Bad JSON: %s", err))
		http.Error(w, "Bad JSON.", 400)
		return
	}

	// retrieve the poll id belonging to the option id
	pollId, err := srv.db.PollIdForOptionId(vote.OptionId)
	if err != nil {
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Unknown option id: %d.",
			pollId))
		http.Error(w, "Unknown option id.", 400)
		return
	}

	// make sure the user is allowed to vote
	if !srv.hasPollAccess(usr.Id, pollId) {
		srv.logger.Log("POST/VOTE", "User has no access rights to the poll.")
		http.Error(w, "Illegal operation.", 400)
		return
	}

	// start a transaction
	transaction, err := srv.db.Begin()
	if err != nil {
		transaction.Rollback()
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err))
		http.Error(w, "Database error", 500)
		return
	}

	// remove all existing votes of the user
	err = database.DelVotesForUserTx(usr.Id, pollId, transaction)
	if err != nil {
		transaction.Rollback()
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err))
		http.Error(w, "Database error", 500)
		return
	}

	// insert the vote into the database
	vote.UserId = usr.Id
	vote.PollId = pollId
	err = database.AddVoteTx(&vote, transaction)
	if err != nil {
		transaction.Rollback()
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err))
		http.Error(w, "Database error", 500)
		return
	}

	// update the poll last updated
	err = database.UpdatePollLastUpdatedTx(pollId, time.Now().Unix(),
		transaction)
	if err != nil {
		transaction.Rollback()
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err))
		http.Error(w, "Database error", 500)
		return
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		srv.logger.Log("POST/VOTE", fmt.Sprintf("Database error: %s.", err))
		http.Error(w, "Database error", 500)
		return
	}

	// send the response
	vote.UserId = usr.Id
	responseBody, err := json.MarshalIndent(vote, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		srv.logger.Log("POST/VOTE",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
		http.Error(w, "Marshalling error.", 500)
		return
	}

}
