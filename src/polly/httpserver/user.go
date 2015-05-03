package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly"

	"github.com/julienschmidt/httprouter"
)

const (
	cGetUserTag    = "GET/USER/XX"
	cUpdateUserTag = "PUT/USER"
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

func (srv *HTTPServer) UpdateUser(writer http.ResponseWriter, req *http.Request,
	_ httprouter.Params) {
	var err error

	// authenticate the user
	usr, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cGetUserTag, err, writer, req)
		return
	}

	// decode the given user
	var updatedUsr polly.PrivateUser
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&updatedUsr)
	if err != nil {
		srv.handleBadRequest(cUpdateUserTag, cBadJSONErr, err, writer, req)
		return
	}

	// update user
	err = srv.db.UpdateUser(usr.Id, updatedUsr.DisplayName,
		updatedUsr.DeviceGUID)
	if err != nil {
		srv.handleDatabaseError(cUpdateUserTag, err, writer, req)
		return
	}

}
