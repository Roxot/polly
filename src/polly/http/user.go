package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	cGetUserTag    = "GET/USER/XX"
	cUpdateUserTag = "PUT/USER"
)

type updateUserFields struct {
	DeviceGUID  *string `json:"device_guid"`
	DisplayName *string `json:"display_name"`
}

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
	var fields updateUserFields
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&fields)
	if err != nil {
		srv.handleBadRequest(cUpdateUserTag, cBadJSONErr, err, writer, req)
		return
	}

	// update display name
	if fields.DeviceGUID != nil {
		err = srv.db.UpdateDeviceGUID(usr.Id, *(fields.DeviceGUID))
		if err != nil {
			srv.handleDatabaseError(cUpdateUserTag, err, writer, req)
		}
	}

	// update device GUID
	if fields.DisplayName != nil {
		err = srv.db.UpdateDisplayName(usr.Id, *(fields.DisplayName))
		if err != nil {
			srv.handleDatabaseError(cUpdateUserTag, err, writer, req)
		}
	}

}
