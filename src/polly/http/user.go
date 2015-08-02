package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"polly"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cGetUserBulkTag = "GET/USERS"
	cUpdateUserTag  = "PUT/USER"
)

func (server *sServer) GetUserBulk(writer http.ResponseWriter,
	request *http.Request, params httprouter.Params) {
	var err error

	// authenticate the user
	_, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cGetUserBulkTag, writer, request)
		return
	}

	// retrieve the list of identifiers
	ids := request.URL.Query()[cID]
	if len(ids) > cBulkUserMax {
		server.respondWithError(ERR_ILL_TOO_MANY_IDS, nil, cGetUserBulkTag,
			writer, request)
		return
	}

	// construct the UserBulk object
	userBulkMsg := polly.UserBulkMessage{}
	userBulkMsg.Users = make([]polly.PublicUser, len(ids))
	for idx, idString := range ids {

		// convert the id to an integer
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			server.respondWithError(ERR_BAD_ID, err, cGetUserBulkTag, writer,
				request)
			return
		}

		// retrieve the user object
		user, err := server.db.GetPublicUserByID(id)
		if err != nil {
			server.respondWithError(ERR_BAD_NO_USER, err, cGetUserBulkTag,
				writer, request)
			return
		}

		userBulkMsg.Users[idx] = *user
	}

	// marshall the response
	responseBody, err := json.MarshalIndent(userBulkMsg, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cGetUserBulkTag, writer,
			request)
		return
	}

	// send the response
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cGetUserBulkTag, writer,
			request)
		return
	}

}

func (server *sServer) UpdateUser(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {
	var err error

	// authenticate the user
	user, errCode := server.authenticateRequest(request)
	if errCode != NO_ERR {
		server.respondWithError(errCode, nil, cUpdateUserTag, writer, request)
		return
	}

	// decode the given user
	var updateUserMsg polly.UpdateUserMessage
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&updateUserMsg)
	if err != nil {
		server.respondWithError(ERR_BAD_JSON, err, cUpdateUserTag, writer,
			request)
		return
	}

	// update display name
	if updateUserMsg.DeviceGUID != nil {
		user.DeviceGUID = *(updateUserMsg.DeviceGUID)
		err = server.db.UpdateDeviceGUID(user.ID, user.DeviceGUID)
		if err != nil {
			server.respondWithError(ERR_INT_DB_UPDATE, err, cUpdateUserTag,
				writer, request)
			return
		}
	}

	// update device GUID
	if updateUserMsg.DisplayName != nil {
		user.DisplayName = *(updateUserMsg.DisplayName)
		err = server.db.UpdateDisplayName(user.ID, user.DisplayName)
		if err != nil {
			server.respondWithError(ERR_INT_DB_UPDATE, err, cUpdateUserTag,
				writer, request)
			return
		}
	}

	// create the response body
	responseBody, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		server.respondWithError(ERR_INT_MARSHALL, err, cUpdateUserTag,
			writer, request)
		return
	}

	// send the user a 200 OK with his user info
	err = server.respondWithJSONBody(writer, responseBody)
	if err != nil {
		server.respondWithError(ERR_INT_WRITE, err, cUpdateUserTag, writer,
			request)
		return
	}

}
