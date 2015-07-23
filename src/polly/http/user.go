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

func (server *sServer) GetUser(writer http.ResponseWriter, request *http.Request,
	params httprouter.Params) {

	// authenticate the user
	_, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cGetUserTag, err, writer, request)
		return
	}

	// load the user from the database
	email := params.ByName(cEmail)
	user, err := server.db.PublicUserByEmail(email)
	if err != nil {
		server.handleErr(cGetUserTag, cNoUserErr,
			fmt.Sprintf(cLogFmt, cNoUserErr, email), 400, writer, request)
		return
	}

	// send the response
	responseBody, err := json.MarshalIndent(user, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cGetUserTag, err, writer, request)
		return
	}

}

func (server *sServer) UpdateUser(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {
	var err error

	// authenticate the user
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cGetUserTag, err, writer, request)
		return
	}

	// decode the given user
	var fields updateUserFields
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&fields)
	if err != nil {
		server.handleBadRequest(cUpdateUserTag, cBadJSONErr, err, writer,
			request)
		return
	}

	// update display name
	if fields.DeviceGUID != nil {
		err = server.db.UpdateDeviceGUID(user.ID, *(fields.DeviceGUID))
		if err != nil {
			server.handleDatabaseError(cUpdateUserTag, err, writer, request)
		}
	}

	// update device GUID
	if fields.DisplayName != nil {
		err = server.db.UpdateDisplayName(user.ID, *(fields.DisplayName))
		if err != nil {
			server.handleDatabaseError(cUpdateUserTag, err, writer, request)
		}
	}

}
