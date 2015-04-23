package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

const (
	cRegisterTag       = "POST/REGISTER"
	cVerifyRegisterTag = "POST/REGISTER/VERIFY"
)

func (srv *HTTPServer) Register(writer http.ResponseWriter, req *http.Request,
	_ httprouter.Params) {

	// validate the provided email address
	email := req.PostFormValue(cEmail)
	if !isValidEmail(email) {
		srv.handleErr(cRegisterTag, cBadEmailErr,
			fmt.Sprintf(cLogFmt, cBadEmailErr, email), 400, writer, req)
		return
	}

	// create a new verification token
	verTkn := polly.VerToken{}
	verTkn.Email = email
	verTkn.VerificationToken = "VERIFY"

	// remove existing verification tokens
	srv.db.DelVerTokensByEmail(email)

	// add the verification token to the database
	err := srv.db.AddVerToken(&verTkn)
	if err != nil {
		srv.handleDatabaseError(cRegisterTag, err, writer, req)
		return
	}

}

func (srv *HTTPServer) VerifyRegister(writer http.ResponseWriter,
	req *http.Request, _ httprouter.Params) {

	retrVerTkn := req.PostFormValue(cVerToken)
	email := req.PostFormValue(cEmail)
	dbVerTkn, err := srv.db.VerTokenByEmail(email)
	if err != nil || dbVerTkn.VerificationToken != retrVerTkn {
		srv.handleErr(cVerifyRegisterTag, cNotRegOrBadTknErr,
			cNotRegOrBadTknErr, 400, writer, req)
		return
	}

	dvcTypeStr := req.PostFormValue(cDeviceType)
	dvcType, err := strconv.Atoi(dvcTypeStr)
	if err != nil || (dvcType != polly.DEVICE_TYPE_AD &&
		dvcType != polly.DEVICE_TYPE_IP) {

		srv.handleErr(cVerifyRegisterTag, cBadDvcTypeErr,
			fmt.Sprintf(cLogFmt, cBadDvcTypeErr, dvcTypeStr), 400, writer, req)
		return
	}

	// device GUID may be empty
	dvcGUID := req.PostFormValue(cDeviceGUID) // TODO validate

	dspName := req.PostFormValue(cDisplayName)
	srv.db.DelVerTokensByEmail(email)
	usr, err := srv.db.UserByEmail(email) // TODO userExists
	if err == nil {

		/* We're dealing with an already existing user */
		responseBody, err := json.MarshalIndent(usr, "", "\t")
		_, err = writer.Write(responseBody)
		if err != nil {
			srv.handleWritingError(cVerifyRegisterTag, err, writer, req)
			return
		}

	} else {

		/* We're dealing with a new user. */
		usr := polly.PrivateUser{}
		usr.Email = email
		usr.Token = uuid.NewV4().String()
		usr.DisplayName = dspName
		usr.DeviceType = dvcType
		usr.DeviceGUID = dvcGUID
		err = srv.db.AddUser(&usr)
		if err != nil {
			srv.handleDatabaseError(cVerifyRegisterTag, err, writer, req)
			return
		}

		responseBody, err := json.MarshalIndent(usr, "", "\t")
		_, err = writer.Write(responseBody)
		if err != nil {
			srv.handleWritingError(cVerifyRegisterTag, err, writer, req)
			return
		}
	}
}
