package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	cGetUserTag = "GET/USER/XX"
)

func (srv *HTTPServer) GetUser(writer http.ResponseWriter, req *http.Request,
	params httprouter.Params) {

	// authenticate the user
	_, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cGetUserTag, err, writer, req)
		return
	}

	// load the user from the database
	email := params.ByName(cEmail)
	usr, err := srv.db.PublicUserByEmail(email)
	if err != nil {
		srv.handleErr(cGetUserTag, cNoUserErr,
			fmt.Sprintf(cLogFmt, cNoUserErr, email), 400, writer, req)
		return
	}

	// send the response
	responseBody, err := json.MarshalIndent(usr, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		srv.handleWritingError(cGetUserTag, err, writer, req)
		return
	}

}
