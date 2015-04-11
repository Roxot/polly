package httpserver

import (
	"fmt"
	"net/http"
	"polly"
)

func (srv *HTTPServer) authenticateRequest(req *http.Request) (
	*polly.PrivateUser, error) {

	phoneNo, tkn, ok := req.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("No authentication provided.")
	}

	usr, err := srv.db.UserByPhoneNumber(phoneNo)
	if err != nil {
		return nil, fmt.Errorf("Unknown user: %s.", phoneNo)
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
			"Somehow existsParticipant returned an error")
		return false
	}

	return exists
}
