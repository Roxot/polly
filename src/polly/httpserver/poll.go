package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly/database"
	"strconv"

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
		srv.logger.Log("POST/POLL", fmt.Sprintf("Unknown user: %s.", phoneNumber))
		http.Error(w, "Unknown user.", 400)
		return
	}

	if user.Token != token {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Bad token: %s doesn't match %s.", token, user.Token))
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

	// if phoneNumber != pollData.Creator.PhoneNumber {
	// 	srv.logger.Log("POST/POLL",
	// 		fmt.Sprintf("Illegal operation: %s and %s don't match.",
	// 			phoneNumber, pollData.Creator.PhoneNumber))
	// 	http.Error(w, "Illegal operation.", 400)
	// 	return
	// }

	pollData.Creator.Id = user.Id

	// Validate poll (non-empty title, no votes (?), etc.)

	err = srv.db.InsertPollData(&pollData)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Database error: %s.", err))
		http.Error(w, "Database error.", 500)
		return
	}
}

func (srv *HTTPServer) GetPoll(w http.ResponseWriter, r *http.Request,
	p httprouter.Params) {

	phoneNumber, token, ok := r.BasicAuth()
	if !ok {
		srv.logger.Log("GET/POLL/XX", "No authentication provided.")
		http.Error(w, "No authentication provided.", 400)
		return
	}

	user, err := srv.db.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Unknown user: %s.", phoneNumber))
		http.Error(w, "Unknown user.", 400)
		return
	}

	if user.Token != token {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Bad token: %s doesn't match %s.", token, user.Token))
		http.Error(w, "Bad token.", 400)
		return
	}

	idString := p.ByName("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Bad id: %s", idString))
		http.Error(w, "Bad id.", 400)
	}

	// TODO Maybe check whether legal more efficiently here

	pollData, err := srv.db.RetrievePollData(id)
	if err != nil {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("No poll with id %s: %s",
			idString, err))
		http.Error(w, "No such poll.", 400)
	}

	// TODO not only creator, also participants
	if pollData.Creator.PhoneNumber != phoneNumber {
		srv.logger.Log("GET/POLL/XX",
			fmt.Sprintf("Illegal operation: retrieving poll from %s while being %s",
				pollData.Creator.PhoneNumber, phoneNumber))
		http.Error(w, "Illegal operation.", 400)
	}

	responseBody, err := json.MarshalIndent(pollData, "", "\t")
	_, err = w.Write(responseBody)
}
