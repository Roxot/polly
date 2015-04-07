package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"polly/database"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func (srv *HTTPServer) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	phoneNumber := r.PostFormValue(cPhoneNumber)
	if !isValidPhoneNumber(phoneNumber) {
		srv.logger.Log("POST/REGISTER", fmt.Sprintf("Bad phone number: %s",
			phoneNumber))
		http.Error(w, "Bad phone number.", 400)
	} else {
		vt := database.VerificationToken{}
		vt.PhoneNumber = phoneNumber
		vt.VerificationToken = "VERIFY"
		srv.db.DeleteVerificationTokensByPhoneNumber(&vt)
		err := srv.db.AddVerificationToken(&vt)
		if err != nil {
			srv.logger.Log("POST/REGISTER", fmt.Sprintf("DATABASE ERROR: %s",
				err))
			http.Error(w, "Database error.", 500)
		}
	}
}

func (srv *HTTPServer) VerifyRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	verificationToken := r.PostFormValue(cVerToken)
	phoneNumber := r.PostFormValue(cPhoneNumber)
	vt, err := srv.db.FindVerificationTokenByPhoneNumber(phoneNumber)
	if err != nil || vt.VerificationToken != verificationToken {
		srv.logger.Log("POST/REGISTER/VERIFY",
			fmt.Sprintf("Not registered or bad verification token: %s - %s",
				phoneNumber, verificationToken))
		http.Error(w, "Not registered or bad verification token.", 400)
		return
	}

	deviceType, err := strconv.Atoi(r.PostFormValue(cDeviceType))
	if err != nil || deviceType < 0 || deviceType > 1 {
		srv.logger.Log("POST/REGISTER/VERIFY",
			fmt.Sprintf("Bad device type: %s",
				r.PostFormValue("device_type")))
		http.Error(w, "Bad device type.", 400)
		return
	}

	displayName := r.PostFormValue(cDisplayName)

	srv.db.DeleteVerificationTokensByPhoneNumber(&vt)
	user, err := srv.db.FindUserByPhoneNumber(phoneNumber)
	if err == nil {

		/* We're dealing with an already existing user */
		uwt := UserToUserWithToken(user)

		responseBody, err := json.MarshalIndent(uwt, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			srv.logger.Log("POST/REGISTER/VERIFY",
				fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
			http.Error(w, "Marshalling error.", 500)
		}
	} else {

		/* We're dealing with a new user. */
		user := database.User{}
		user.PhoneNumber = phoneNumber
		user.Token = uuid.NewV4().String()
		user.DisplayName = displayName
		user.DeviceType = deviceType
		err = srv.db.AddUser(&user)
		if err != nil {
			srv.logger.Log("POST/REGISTER/VERIFY",
				fmt.Sprintf("DATABASE ERROR: %s\n", err))
			http.Error(w, "Database error.", 500)
			return
		}

		uwt := UserToUserWithToken(user)

		responseBody, err := json.MarshalIndent(uwt, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			srv.logger.Log("POST/REGISTER/VERIFY",
				fmt.Sprintf("MARSHALLING ERROR: %s\n", err))
			http.Error(w, "Marshalling error.", 500)
		}
	}
}
