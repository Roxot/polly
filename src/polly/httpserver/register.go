package httpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"polly"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func (srv *HTTPServer) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	phoneNo := r.PostFormValue(cPhoneNumber)
	if !isValidPhoneNumber(phoneNo) {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/REGISTER", fmt.Sprintf("Bad phone number: %s",
			phoneNo), h)
		http.Error(w, "Bad phone number.", 400)
		return
	} else {
		verTkn := polly.VerToken{}
		verTkn.PhoneNumber = phoneNo
		verTkn.VerificationToken = "VERIFY"
		srv.db.DelVerTokensByPhoneNumber(phoneNo)
		err := srv.db.AddVerToken(&verTkn)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/REGISTER", fmt.Sprintf("DATABASE ERROR: %s",
				err), h)
			http.Error(w, "Database error.", 500)
			return
		}
	}
}

func (srv *HTTPServer) VerifyRegister(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	retrVerTkn := r.PostFormValue(cVerToken)
	phoneNo := r.PostFormValue(cPhoneNumber)
	dbVerTkn, err := srv.db.VerTokenByPhoneNumber(phoneNo)
	if err != nil || dbVerTkn.VerificationToken != retrVerTkn {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/REGISTER/VERIFY",
			fmt.Sprintf("Not registered or bad verification token: %s - %s",
				phoneNo, retrVerTkn), h)
		http.Error(w, "Not registered or bad verification token.", 400)
		return
	}

	dvcType, err := strconv.Atoi(r.PostFormValue(cDeviceType))
	if err != nil || dvcType < 0 || dvcType > 1 {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/REGISTER/VERIFY",
			fmt.Sprintf("Bad device type: %s",
				r.PostFormValue("device_type")), h)
		http.Error(w, "Bad device type.", 400)
		return
	}

	dspName := r.PostFormValue(cDisplayName)
	srv.db.DelVerTokensByPhoneNumber(phoneNo)
	usr, err := srv.db.UserByPhoneNumber(phoneNo)
	if err == nil {

		/* We're dealing with an already existing user */
		responseBody, err := json.MarshalIndent(usr, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/REGISTER/VERIFY",
				fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
			http.Error(w, "Marshalling error.", 500)
		}
	} else {

		/* We're dealing with a new user. */
		usr := polly.PrivateUser{}
		usr.PhoneNumber = phoneNo
		usr.Token = uuid.NewV4().String()
		usr.DisplayName = dspName
		usr.DeviceType = dvcType
		usr.DeviceGUID = ""
		err = srv.db.AddUser(&usr)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/REGISTER/VERIFY",
				fmt.Sprintf("DATABASE ERROR: %s\n", err), h)
			http.Error(w, "Database error.", 500)
			return
		}

		responseBody, err := json.MarshalIndent(usr, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			h, _, _ := net.SplitHostPort(r.RemoteAddr)
			srv.logger.Log("POST/REGISTER/VERIFY",
				fmt.Sprintf("MARSHALLING ERROR: %s\n", err), h)
			http.Error(w, "Marshalling error.", 500)
		}
	}
}
