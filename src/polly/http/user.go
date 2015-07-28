package http

import (
	"encoding/json"
	"net/http"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cGetUserTag    = "GET/USER/XX"
	cUpdateUserTag = "PUT/USER"
)

type sUpdateUserFields struct {
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
	// email := params.ByName(cEmail)
	// user, err := server.db.GetPublicUserByEmail(email)
	// if err != nil {
	// 	server.handleErr(cGetUserTag, cNoUserErr,
	// 		fmt.Sprintf(cLogFmt, cNoUserErr, email), 400, writer, request)
	// 	return
	// }

	// // send the response
	// responseBody, err := json.MarshalIndent(user, "", "\t")
	// _, err = writer.Write(responseBody)
	// if err != nil {
	// 	server.handleWritingError(cGetUserTag, err, writer, request)
	// 	return
	// }

	// TODO fix
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
	var fields sUpdateUserFields
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&fields)
	if err != nil {
		server.handleBadRequest(cUpdateUserTag, cBadJSONErr, err, writer,
			request)
		return
	}

	// update display name
	if fields.DeviceGUID != nil {
		user.DeviceGUID = *(fields.DeviceGUID)
		err = server.db.UpdateDeviceGUID(user.ID, user.DeviceGUID)
		if err != nil {
			server.handleDatabaseError(cUpdateUserTag, err, writer, request)
		}
	}

	// update device GUID
	if fields.DisplayName != nil {
		user.DisplayName = *(fields.DisplayName)
		err = server.db.UpdateDisplayName(user.ID, user.DisplayName)
		if err != nil {
			server.handleDatabaseError(cUpdateUserTag, err, writer, request)
		}
	}

	// create the response body
	responseBody, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		server.handleMarshallingError(cRegisterTag, err, writer, request)
		return
	}

	// send the user a 200 OK with his user info
	SetJSONContentType(writer)
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cRegisterTag, err, writer, request)
		return
	}

}
