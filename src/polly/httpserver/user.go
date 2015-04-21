package httpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srv *HTTPServer) GetUser(w http.ResponseWriter, r *http.Request,
	p httprouter.Params) {

	// authenticate the user
	_, err := srv.authenticateRequest(r)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/USER/XX", fmt.Sprintf("Authentication error: %s",
			err), h)
		w.Header().Set("WWW-authenticate", "Basic")
		http.Error(w, "Authentication error", 401)
		return
	}

	// retrieve the email parameter
	email := p.ByName(cEmail)
	if len(email) == 0 {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/USER/XX", "Empty email.", h)
		http.Error(w, "Bad email.", 400)
		return
	}

	// load the user from the database
	usr, err := srv.db.PublicUserByEmail(email)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/USER/XX", fmt.Sprintf("User not found: %s.",
			email), h)
		http.Error(w, "Unkown user.", 400)
		return
	}

	// send the response
	responseBody, err := json.MarshalIndent(usr, "", "\t")
	_, err = w.Write(responseBody)
	if err != nil {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("GET/USER/XX",
			fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
		http.Error(w, "Marshalling error.", 500)
		return
	}

}
