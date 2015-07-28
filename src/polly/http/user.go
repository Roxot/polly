package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"polly"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cGetUserBulkTag = "GET/USERS"
	cUpdateUserTag  = "PUT/USER"
)

type sUpdateUserFields struct {
	DeviceGUID  *string `json:"device_guid"`
	DisplayName *string `json:"display_name"`
}

type UserBulk struct { // TODO private or model
	Users []polly.PublicUser `json:"users"`
}

func (server *sServer) GetUserBulk(writer http.ResponseWriter,
	request *http.Request, params httprouter.Params) {

	// authenticate the user
	_, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cGetUserBulkTag, err, writer, request)
		return
	}

	// retrieve the list of identifiers
	ids := request.URL.Query()[cID]
	if len(ids) > cBulkUserMax {
		server.handleErr(cGetUserBulkTag, cIDListLengthErr,
			fmt.Sprintf("%s: %d", cIDListLengthErr, len(ids)), 400, writer,
			request)
		return
	}

	// construct the UserBulk object
	userBulk := UserBulk{}
	userBulk.Users = make([]polly.PublicUser, len(ids))
	for idx, idString := range ids {

		// convert the id to an integer
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			server.handleBadRequest(cGetUserBulkTag, cBadIDErr, err, writer,
				request)
			return
		}

		// retrieve the user object
		user, err := server.db.GetPublicUserByID(id)
		if err != nil {
			server.handleErr(cGetUserBulkTag, cNoUserErr, cNoUserErr, 400,
				writer, request)
			return
		}

		userBulk.Users[idx] = *user
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(userBulk, "", "\t")
	if err != nil {
		server.handleMarshallingError(cGetUserBulkTag, err, writer, request)
		return
	}

	// send the response
	SetJSONContentType(writer)
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cGetUserBulkTag, err, writer, request)
		return
	}

}

func (server *sServer) UpdateUser(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {
	var err error

	// authenticate the user
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cUpdateUserTag, err, writer, request)
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
