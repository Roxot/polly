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
	phoneNumber := r.PostFormValue("phone_number")
	if !isValidPhoneNumber(phoneNumber) {
		srv.logger.Log(fmt.Sprintf("[POST/REGISTER] invalid phonenumber %s\n", phoneNumber))
		http.Error(w, "Invalid phonenumber.", 400)
	} else {
		vt := database.VerificationToken{}
		vt.PhoneNumber = phoneNumber
		vt.VerificationToken = "VERIFY"
		srv.db.DeleteVerificationTokensByPhoneNumber(&vt)
		err := srv.db.AddVerificationToken(&vt)
		if err != nil {
			srv.logger.Log(fmt.Sprintf("[POST/REGISTER] DATABASE ERROR %s\n", err))
			http.Error(w, "Database error.", 500)
		}
	}
}

func (srv *HTTPServer) VerifyRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	verificationToken := r.PostFormValue("verification_token")
	phoneNumber := r.PostFormValue("phone_number")
	vt, err := srv.db.FindVerificationTokenByPhoneNumber(phoneNumber)
	if err != nil || vt.VerificationToken != verificationToken {
		srv.logger.Log(
			fmt.Sprintf("[POST/REGISTER/VERIFY] not registered / bad token %s - %s\n",
				phoneNumber, verificationToken))
		http.Error(w, "Not registered / bad token.", 400)
		return
	}

	deviceType, err := strconv.Atoi(r.PostFormValue("device_type"))
	if err != nil || deviceType < 0 || deviceType > 1 {
		srv.logger.Log(fmt.Sprintf("[POST/REGISTER/VERIFY] bad device type %s\n",
			r.PostFormValue("device_type")))
		http.Error(w, "Invalid device type.", 400)
		return
	}

	displayName := r.PostFormValue("display_name")

	srv.db.DeleteVerificationTokensByPhoneNumber(&vt)
	user, err := srv.db.FindUserByPhoneNumber(phoneNumber)
	if err == nil {

		/* We're dealing with an already existing user */
		uwt := UserToUserWithToken(user)

		responseBody, err := json.MarshalIndent(uwt, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			srv.logger.Log(fmt.Sprintf("[POST/REGISTER/VERIFY] MARSHALLING ERROR %s\n", err))
			http.Error(w, "Marshalling error.", 500)
		}
	} else {
		// new user
		user := database.User{}
		user.PhoneNumber = phoneNumber
		user.Token = uuid.NewV4().String()
		user.DisplayName = displayName
		user.DeviceType = deviceType
		err = srv.db.AddUser(&user)
		if err != nil {
			srv.logger.Log(fmt.Sprintf("[POST/REGISTER/VERIFY] DATABASE ERROR %s\n", err))
			http.Error(w, "Database error.", 500)
			return
		}

		uwt := UserToUserWithToken(user)

		responseBody, err := json.MarshalIndent(uwt, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			srv.logger.Log(fmt.Sprintf("[POST/REGISTER/VERIFY] MARSHALLING ERROR %s\n", err))
			http.Error(w, "Marshalling error.", 500)
		}
	}
}
