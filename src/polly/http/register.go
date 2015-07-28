package http

import (
	"encoding/json"
	"net/http"
	"polly"

	"polly/internal/github.com/julienschmidt/httprouter"
	"polly/internal/github.com/satori/go.uuid"
)

const (
	cRegisterTag                  = "POST/REGISTER"
	cVerifyRegisterTag            = "POST/REGISTER/VERIFY"
	cServiceProviderHeader        = "X-Auth-Service-Provider"
	cClientAuthorizationHeader    = "X-Verify-Credentials-Authorization"
	cOAuthEchoAuthorizationHeader = "Authorization" // TODO correct names?
)

type tTwitterResponse struct {
	IDString    string
	ID          int64
	CreatedAt   string
	PhoneNumber string
	AccessToken tTwitterAccessToken
}

type tTwitterAccessToken struct {
	Secret string
	Token  string
}

// TODO security checks headers
func (server *sServer) Register(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {
	var decoder *json.Decoder

	// retrieve headers for the OAuth Echo
	provider := request.Header.Get(cServiceProviderHeader)
	authorization := request.Header.Get(cClientAuthorizationHeader)
	if len(provider) == 0 || len(authorization) == 0 {
		// TODO
		server.handleErr(cVerifyRegisterTag, cAuthErr, cAuthErr, 403, writer,
			request)
		return
	}

	// create an HTTP Client and the OAuth Echo request
	client := http.Client{}
	oauthEcho, err := http.NewRequest("GET", provider, nil)
	if err != nil {
		// TODO wrong error function
		server.handleDatabaseError(cVerifyRegisterTag, err, writer, request)
		return
	}

	// add the authorization header to the OAuth Echo request
	oauthEcho.Header.Add(cOAuthEchoAuthorizationHeader, authorization)

	// verify the provided credentials with Twitter's server
	response, err := client.Do(oauthEcho)
	if err != nil {
		// TODO wrong error function
		server.handleDatabaseError(cVerifyRegisterTag, err, writer, request)
		return
	} else if response.StatusCode != 200 { // TODO status code constants somewhere..
		// TODO
		server.handleErr(cVerifyRegisterTag, cAuthErr, cAuthErr, 403, writer,
			request)
		return
	}

	// decode Twitter's response
	var twitterResponse tTwitterResponse
	decoder = json.NewDecoder(response.Body)
	err = decoder.Decode(&twitterResponse)
	if err != nil {
		// TODO
		server.handleBadRequest(cVerifyRegisterTag, cAuthErr, err, writer,
			request)
		return
	}

	// decode the sent user info
	var user polly.PrivateUser
	decoder = json.NewDecoder(request.Body)
	err = decoder.Decode(&user)
	if err != nil {
		// TODO
		server.handleBadRequest(cVerifyRegisterTag, cBadJSONErr, err, writer,
			request)
		return
	}

	// check that the device type is a correct one
	if !isValidDeviceType(user.DeviceType) {
		server.handleErr(cVerifyRegisterTag, cBadDvcTypeErr,
			cBadDvcTypeErr, 400, writer, request) // TODO log device type?
		return
	}

	existentUser, err := server.db.GetUserByID(twitterResponse.ID)
	if err == nil {

		// we're dealing with an already existing user, generate a new token
		// and update his device type (TODO is this needed?)
		existentUser.DeviceType = user.DeviceType
		existentUser.Token = uuid.NewV4().String()
		err = server.db.UpdateToken(existentUser.ID, existentUser.Token)
		if err != nil {
			server.handleDatabaseError(cVerifyRegisterTag, err, writer, request)
			return
		}

		// send the user a 200 OK with his new user info
		SetJSONContentType(writer)
		responseBody, err := json.MarshalIndent(existentUser, "", "\t")
		_, err = writer.Write(responseBody)
		if err != nil {
			server.handleWritingError(cVerifyRegisterTag, err, writer, request)
			return
		}

	} else {

		// we're dealing with a new user
		user.Token = uuid.NewV4().String()
		user.ID = twitterResponse.ID
		err = server.db.AddUser(&user)
		if err != nil {
			server.handleDatabaseError(cVerifyRegisterTag, err, writer, request)
			return
		}

		// send the user a 200 OK with his user info
		SetJSONContentType(writer)
		responseBody, err := json.MarshalIndent(user, "", "\t")
		_, err = writer.Write(responseBody)
		if err != nil {
			server.handleWritingError(cVerifyRegisterTag, err, writer, request)
			return
		}
	}

}

// TODO move
func SetJSONContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
