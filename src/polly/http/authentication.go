package httpserver

import (
	"fmt"
	"net/http"
	"polly"
)

func (srv *HTTPServer) authenticateRequest(req *http.Request) (
	*polly.PrivateUser, error) {

	email, tkn, ok := req.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("No authentication provided.")
	}

	usr, err := srv.db.UserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("Unknown user: %s.", email)
	}

	if usr.Token != tkn {
		return nil, fmt.Errorf("Bad token.")
	}

	return usr, nil
}

func (srv *HTTPServer) hasPollAccess(usrId, pollId int) bool {
	exists, err := srv.db.ExistsParticipant(usrId, pollId)
	if err != nil {
		srv.logger.Log("hasPollAccess",
			"Somehow existsParticipant returned an error", "hasPollAccess")
		return false
	}

	return exists
}
