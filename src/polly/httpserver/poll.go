package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly/database"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (srv *HTTPServer) PostPoll(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	err := srv.authenticateUser(r)
	if err != nil {
		srv.logger.Log("POST/POLL", fmt.Sprintf("Authentication error: %s",
			err))
		http.Error(w, "Authentication error", 400)
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
	pollData.Creator.PhoneNumber = user.PhoneNumber
	pollData.Creator.DisplayName = user.DisplayName
	pollData.MetaData.CreationDate = time.Now().Unix()

	// Validate poll (non-empty title, no votes (?), etc.)

	err, isInternalErr := srv.db.InsertPollData(&pollData)
	if err != nil {
		if isInternalErr {
			srv.logger.Log("POST/POLL", fmt.Sprintf("Database error: %s.", err))
			http.Error(w, "Database error.", 500)
			return
		} else {
			srv.logger.Log("POST/POLL", fmt.Sprintf("Bad poll: %s.", err))
			http.Error(w, "Bad poll.", 400)
			return
		}
	}

	responseBody, err := json.MarshalIndent(pollData, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		srv.logger.Log("POST/POLL",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
		http.Error(w, "Marshalling error.", 500)
	}

}

func (srv *HTTPServer) GetPoll(w http.ResponseWriter, r *http.Request,
	p httprouter.Params) {

	err := srv.authenticateUser(r)
	if err != nil {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Authentication error: %s",
			err))
		http.Error(w, "Authentication error", 400)
		return
	}

	idString := p.ByName(cId)
	id, err := strconv.Atoi(idString)
	if err != nil {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("Bad id: %s", idString))
		http.Error(w, "Bad id.", 400)
		return
	}

	// TODO Maybe check whether legal more efficiently here

	pollData, err := srv.db.RetrievePollData(id)
	if err != nil {
		srv.logger.Log("GET/POLL/XX", fmt.Sprintf("No poll with id %s: %s",
			idString, err))
		http.Error(w, "No such poll.", 400)
		return
	}

	// TODO not only creator, also participants
	if pollData.Creator.PhoneNumber != phoneNumber {
		srv.logger.Log("GET/POLL/XX",
			fmt.Sprintf("Illegal operation: retrieving poll from %s while being %s",
				pollData.Creator.PhoneNumber, phoneNumber))
		http.Error(w, "Illegal operation.", 400)
		return
	}

	responseBody, err := json.MarshalIndent(pollData, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		srv.logger.Log("GET/POLL/XX",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
		http.Error(w, "Marshalling error.", 500)
	}

}
