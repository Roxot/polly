package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PollInfo struct {
	PollId      int `json:"poll_id"`
	LastUpdated int `json:"last_updated"`
}

type PollList struct {
	PollInfos []PollInfo `json:"polls"`
}

func (srv *HTTPServer) ListUserPolls(w http.ResponseWriter, r *http.Request,
	p httprouter.Params) {

	err := srv.authenticateUser(r)
	if err != nil {
		srv.logger.Log("USER/POLLS", fmt.Sprintf("Authentication error: %s",
			err))
		http.Error(w, "Authentication error", 400)
		return
	}

	polls, err := srv.db.FindPollsByUserId(user.Id)
	if err != nil {
		log.Fatal(err)
	}

	pollList := PollList{}
	pollList.PollInfos = make([]PollInfo, len(polls))
	for index, poll := range polls {
		pollList.PollInfos[index].PollId = poll.Id
		pollList.PollInfos[index].LastUpdated = 0
	}

	responseBody, err := json.MarshalIndent(pollList, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		srv.logger.Log("USER/POLLS",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
		http.Error(w, "Marshalling error.", 500)
	}
}
