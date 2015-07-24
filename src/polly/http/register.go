package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly"
	"strconv"

	"polly/internal/github.com/julienschmidt/httprouter"
	"polly/internal/github.com/satori/go.uuid"
)

const (
	cRegisterTag       = "POST/REGISTER"
	cVerifyRegisterTag = "POST/REGISTER/VERIFY"
)

func (server *sServer) Register(writer http.ResponseWriter, request *http.Request,
	_ httprouter.Params) {

	// validate the provided email address
	email := request.PostFormValue(cEmail)
	if !isValidEmail(email) {
		server.handleErr(cRegisterTag, cBadEmailErr,
			fmt.Sprintf(cLogFmt, cBadEmailErr, email), 400, writer, request)
		return
	}

	// create a new verification token
	verToken := polly.VerToken{}
	verToken.Email = email
	verToken.Token = "VERIFY"

	// remove existing verification tokens
	server.db.DeleteVerTokensByEmail(email)

	// add the verification token to the database
	err := server.db.AddVerToken(&verToken)
	if err != nil {
		server.handleDatabaseError(cRegisterTag, err, writer, request)
		return
	}

}

func (server *sServer) VerifyRegister(writer http.ResponseWriter,
	request *http.Request, _ httprouter.Params) {

	retrVerToken := request.PostFormValue(cVerToken)
	email := request.PostFormValue(cEmail)
	dbVerToken, err := server.db.GetVerTokenByEmail(email)
	if err != nil || dbVerToken.Token != retrVerToken {
		server.handleErr(cVerifyRegisterTag, cNotRegOrBadTknErr,
			cNotRegOrBadTknErr, 400, writer, request)
		return
	}

	deviceTypeStr := request.PostFormValue(cDeviceType)
	deviceType, err := strconv.Atoi(deviceTypeStr)
	if err != nil || (deviceType != polly.DEVICE_TYPE_AD &&
		deviceType != polly.DEVICE_TYPE_IP) {

		server.handleErr(cVerifyRegisterTag, cBadDvcTypeErr,
			fmt.Sprintf(cLogFmt, cBadDvcTypeErr, deviceTypeStr), 400, writer,
			request)
		return
	}

	// device GUID may be empty
	deviceGUID := request.PostFormValue(cDeviceGUID) // TODO validate

	displayName := request.PostFormValue(cDisplayName)
	server.db.DeleteVerTokensByEmail(email)
	user, err := server.db.GetUserByEmail(email) // TODO userExists
	if err == nil {

		/* We're dealing with an already existing user */
		responseBody, err := json.MarshalIndent(user, "", "\t")
		_, err = writer.Write(responseBody)
		if err != nil {
			server.handleWritingError(cVerifyRegisterTag, err, writer, request)
			return
		}

	} else {

		/* We're dealing with a new user. */
		user := polly.PrivateUser{}
		user.Email = email
		user.Token = uuid.NewV4().String()
		user.DisplayName = displayName
		user.DeviceType = deviceType
		user.DeviceGUID = deviceGUID
		err = server.db.AddUser(&user)
		if err != nil {
			server.handleDatabaseError(cVerifyRegisterTag, err, writer, request)
			return
		}

		responseBody, err := json.MarshalIndent(user, "", "\t")
		_, err = writer.Write(responseBody)
		if err != nil {
			server.handleWritingError(cVerifyRegisterTag, err, writer, request)
			return
		}
	}
}
