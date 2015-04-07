package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PollBulk struct {
	Polls []database.PollData `json:"polls"`
}

func (srv *HTTPServer) GetPollBulk(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	err := srv.authenticateUser(r)
	if err != nil {
		srv.logger.Log("GET/POLL", fmt.Sprintf("Authentication error: %s",
			err))
		http.Error(w, "Authentication error", 400)
		return
	}

	ids := r.URL.Query()[cId]
	if len(ids) > cBulkPollMax {
		srv.logger.Log("GET/POLL",
			fmt.Sprintf("Id list longer than limit: %d > %d", len(ids), cBulkPollMax))
		http.Error(w, "Id list longer than limit", 400)
		return
	}

	pollBulk := PollBulk{}
	pollBulk.Polls = make([]database.PollData, len(ids))
	for index, idString := range ids {

		id, err := strconv.Atoi(idString)
		if err != nil {
			srv.logger.Log("GET/POLL", fmt.Sprintf("Bad id: %s", idString))
			http.Error(w, "Bad id.", 400)
			return
		}

		pollData, err := srv.db.RetrievePollData(id)
		if err != nil {
			srv.logger.Log("GET/POLL", fmt.Sprintf("No poll with id %s: %s",
				idString, err))
			http.Error(w, "No such poll.", 400)
			return
		}

		// TODO not only creator, also participants
		if pollData.Creator.PhoneNumber != phoneNumber {
			srv.logger.Log("GET/POLL",
				fmt.Sprintf("Illegal operation: retrieving poll from %s while being %s",
					pollData.Creator.PhoneNumber, phoneNumber))
			http.Error(w, "Illegal operation.", 400)
			return
		}

		pollBulk.Polls[index] = pollData
	}

	responseBody, err := json.MarshalIndent(pollBulk, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		srv.logger.Log("GET/POLL",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
		http.Error(w, "Marshalling error.", 500)
	}
}
