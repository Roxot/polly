package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly/database"

	"github.com/julienschmidt/httprouter"
)

func (srv *HTTPServer) PostPoll(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	phoneNumber, token, ok := r.BasicAuth()
	if !ok {
		srv.logger.Log("POST/POLL", "No authentication provided.")
		http.Error(w, "No authentication provided.", 400)
		return
	}

	user, err := srv.db.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Unknown user: %s", phoneNumber))
		http.Error(w, "Unknown user.", 400)
		return
	}

	if user.Token != token {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Bad token: %s doesn't match %s", token, user.Token))
		http.Error(w, "Bad token.", 400)
		return
	}

	var pollData database.PollData
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&pollData)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Bad JSON: %s", err))
		http.Error(w, "Bad JSON.", 400)
		return
	}

	if phoneNumber != pollData.Creator.PhoneNumber {
		srv.logger.Log("POST/POLL",
			fmt.Sprintf("Illegal operation: %s and %s don't match",
				phoneNumber, pollData.Creator.PhoneNumber))
		http.Error(w, "Illegal operation.", 400)
		return
	}

	// Validate poll (non-empty title, no votes (?), etc.)

	err = srv.db.InsertPollData(&pollData)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Database error: %s", err))
		http.Error(w, "Database error.", 500)
		return
	}
}
