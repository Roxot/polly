package httpserver

import (
	"fmt"
	"net/http"
)

func (srv *HTTPServer) authenticateUser(req *http.Request) error {

	phoneNumber, token, ok := req.BasicAuth()
	if !ok {
		return fmt.Errorf("No authentication provided.")
	}

	user, err := srv.db.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf("Unknown user: %s.", phoneNumber)
	}

	if user.Token != token {
		return fmt.Errorf("Bad token.")
	}

	return nil
}
