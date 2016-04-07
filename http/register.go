package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/roxot/polly"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

const (
	cRegisterTag                 = "POST/REGISTER"
	cFacebookVerifyTokenEndpoint = "https://graph.facebook.com/me?fields=id&access_token=%s"
	cFacebookVerifyTokenTimeout  = 10 * time.Second
)

type sFacebookResponse struct {
	ID string `json:"id"`
}

func (server *sServer) Register(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	// retrieve the Facebook auth token
	fbToken := request.Header.Get("X-Verify-Credentials-Authorization")
	if len(fbToken) == 0 {
		server.respondWithError(ERR_AUT_NO_FACEBOOK_TOKEN, nil, cRegisterTag,
			writer, request)
		return
	}

	// verify the Facebook auth token
	id, errCode, err := verifyFacebookUser(fbToken)
	if errCode != NO_ERR {
		server.respondWithError(errCode, err, cRegisterTag, writer, request)
		return
	}

	// decode the sent user info
	var user polly.User
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&user)
	if err != nil {
		server.respondWithError(ERR_BAD_JSON, err, cRegisterTag, writer,
			request)
		return
	}

	// check that the device type is a correct one
	if !isValidDeviceType(user.DeviceType) {
		server.respondWithError(ERR_BAD_DEVICE_TYPE, nil, cRegisterTag, writer,
			request)
		return
	}

	// make sure a display name is set
	if len(user.DisplayName) == 0 {
		server.respondWithError(ERR_BAD_NO_DISPLAY_NAME, nil, cRegisterTag,
			writer, request)
		return
	}

	existentUser, err := server.db.GetUserByID(id)
	if err == nil {

		// we're dealing with an already existing user, generate a new token
		// and update his or her device type
		existentUser.DeviceType = user.DeviceType
		// TODO update profile pic && display name
		// TODO stuff is not updated here? bug?
		existentUser.Token = uuid.NewV4().String()
		err = server.db.UpdateToken(existentUser.ID, existentUser.Token)
		if err != nil {
			server.respondWithError(ERR_INT_DB_UPDATE, err, cRegisterTag,
				writer, request)
			return
		}

		// create the response body
		responseBody, err := json.MarshalIndent(existentUser, "", "\t")
		if err != nil {
			server.respondWithError(ERR_INT_MARSHALL, err, cRegisterTag,
				writer, request)
			return
		}

		// send the user a 200 OK with his new user info
		err = server.respondWithJSONBody(writer, responseBody)
		if err != nil {
			server.respondWithError(ERR_INT_WRITE, err, cRegisterTag, writer,
				request)
			return
		}

	} else {

		// we're dealing with a new user
		user.Token = uuid.NewV4().String()
		user.ID = id
		err = server.db.AddUser(&user)
		if err != nil {
			server.respondWithError(ERR_INT_DB_ADD, err, cRegisterTag, writer,
				request)
			return
		}

		// create the response body
		responseBody, err := json.MarshalIndent(user, "", "\t")
		if err != nil {
			server.respondWithError(ERR_INT_MARSHALL, err, cRegisterTag, writer,
				request)
			return
		}

		// send the user a 200 OK with his user info
		err = server.respondWithJSONBody(writer, responseBody)
		if err != nil {
			server.respondWithError(ERR_INT_WRITE, err, cRegisterTag, writer,
				request)
			return
		}
	}

}

func verifyFacebookUser(token string) (int64, int, error) {
	client := http.Client{
		Timeout: cFacebookVerifyTokenTimeout,
	}

	// Request the Facebook user ID using the token
	resp, err := client.Get(fmt.Sprintf(cFacebookVerifyTokenEndpoint, token))
	if err != nil {
		return 0, ERR_INT_DO_HTTP, nil
	} else if resp.StatusCode != http.StatusOK {
		return 0, ERR_AUT_BAD_FACEBOOK_TOKEN, nil
	}

	// Decode the response
	var respBody sFacebookResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respBody)
	if err != nil {
		return 0, ERR_INT_DEMARSHALL, err
	}

	// Parse the Facebook ID as a 64 bit integer
	id, err := strconv.ParseInt(respBody.ID, 10, 64)
	if err != nil {
		return 0, ERR_INT_PARSE_INT, nil
	}

	return id, NO_ERR, nil
}
