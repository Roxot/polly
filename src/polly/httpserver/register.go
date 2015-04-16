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
	email := r.PostFormValue(cEmail)
	if !isValidEmail(email) {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/REGISTER", fmt.Sprintf("Bad email address: %s",
			email), h)
		http.Error(w, "Bad email address.", 400)
		return
	} else {
		verTkn := polly.VerToken{}
		verTkn.Email = email
		verTkn.VerificationToken = "VERIFY"
		srv.db.DelVerTokensByEmail(email)
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
	email := r.PostFormValue(cEmail)
	dbVerTkn, err := srv.db.VerTokenByEmail(email)
	if err != nil || dbVerTkn.VerificationToken != retrVerTkn {
		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/REGISTER/VERIFY",
			fmt.Sprintf("Not registered or bad verification token: %s - %s",
				email, retrVerTkn), h)
		http.Error(w, "Not registered or bad verification token.", 400)
		return
	}

	dvcType, err := strconv.Atoi(r.PostFormValue(cDeviceType))
	if err != nil || (dvcType != polly.DEVICE_TYPE_AD &&
		dvcType != polly.DEVICE_TYPE_IP) {

		h, _, _ := net.SplitHostPort(r.RemoteAddr)
		srv.logger.Log("POST/REGISTER/VERIFY",
			fmt.Sprintf("Bad device type: %s",
				r.PostFormValue("device_type")), h)
		http.Error(w, "Bad device type.", 400)
		return
	}

	dspName := r.PostFormValue(cDisplayName)
	srv.db.DelVerTokensByEmail(email)
	usr, err := srv.db.UserByEmail(email)
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
		usr.Email = email
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
