package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pollydatabase"
	"strconv"
	"unicode"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

var pollyDb pollydatabase.PollyDatabase

func main() {
	var err error
	pollyDb, err = pollydatabase.New()
	checkErr(err)

	//err = pollyDb.DropTablesIfExists()
	//checkErr(err)

	err = pollyDb.CreateTablesIfNotExists()
	checkErr(err)

	/* Insert user. */
	// newUser := pollydatabase.User{}
	// newUser.PhoneNumber = "0612345678"
	// newUser.Token = "P0lly"
	// newUser.DisplayName = "The Wolf of Wall Street"
	// newUser.DeviceType = 0
	// newUser.DeviceGUID = ""
	// err = pollyDb.AddUser(&newUser)
	// checkErr(err)

	router := httprouter.New()
	router.POST("/register", Register)
	router.POST("/register/verify", VerifyRegister)
	// router.GET("/poll/:id", GetPoll)
	// router.POST("/poll/:id/vote", Vote)
	// router.POST("/create/poll", CreatePoll)
	// router.POST("/create/user", CreateUser)
	// router.GET("/user/:id/", GetUser)
	err = http.ListenAndServe(":8080", router)
	checkErr(err)
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	phoneNumber := r.PostFormValue("phone_number")
	if !isValidPhoneNumber(phoneNumber) {
		http.Error(w, "Invalid phonenumber.", 400)
	} else {
		vt := pollydatabase.VerificationToken{}
		vt.PhoneNumber = phoneNumber
		vt.VerificationToken = "VERIFY"
		pollyDb.DeleteVerificationTokensByPhoneNumber(&vt)
		err := pollyDb.AddVerificationToken(&vt)
		if err != nil {
			http.Error(w, "Database error.", 500)
		}
	}
}

func VerifyRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	verificationToken := r.PostFormValue("verification_token")
	phoneNumber := r.PostFormValue("phone_number")
	vt, err := pollyDb.FindVerificationTokenByPhoneNumber(phoneNumber)
	if err != nil || vt.VerificationToken != verificationToken {
		http.Error(w, "Not registered / bad token.", 400)
		log.Println(err)
		log.Printf("%s != %s\n", vt.VerificationToken, verificationToken)
		log.Println("Not registered / bad token.")
		return
	} else {
		pollyDb.DeleteVerificationTokensByPhoneNumber(&vt)
	}

	deviceType, err := strconv.Atoi(r.PostFormValue("device_type"))
	if err != nil || deviceType < 0 || deviceType > 1 {
		log.Println("Invalid device type.")
		http.Error(w, "Invalid device type.", 400)
		return
	}

	displayName := r.PostFormValue("display_name")

	user, err := pollyDb.FindUserByPhoneNumber(phoneNumber)
	if err == nil {

		/* We're dealing with an already existing user */
		uwt := UserToUserWithToken(user)

		responseBody, err := json.MarshalIndent(uwt, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			http.Error(w, "Marshalling error.", 500)
			log.Println(err)
		}
	} else {
		// new user
		user := pollydatabase.User{}
		user.PhoneNumber = phoneNumber
		user.Token = uuid.NewV4().String()
		user.DisplayName = displayName
		user.DeviceType = deviceType
		err = pollyDb.AddUser(&user)
		if err != nil {
			http.Error(w, "Database error.", 500)
			log.Println(err)
			return
		}

		uwt := UserToUserWithToken(user)

		responseBody, err := json.MarshalIndent(uwt, "", "\t")
		_, err = w.Write(responseBody)
		if err != nil {
			http.Error(w, "Marshalling error.", 500)
			log.Println(err)
		}
	}
}

func UserToUserWithToken(user pollydatabase.User) UserWithToken {
	uwt := UserWithToken{}
	uwt.Id = user.Id
	uwt.PhoneNumber = user.PhoneNumber
	uwt.Token = user.Token
	uwt.DisplayName = user.DisplayName
	uwt.DeviceType = user.DeviceType
	return uwt
}

type UserWithToken struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phone_numer"`
	Token       string `json:"token"`
	DisplayName string `json:"display_name"`
	DeviceType  int    `json:"device_type"`
}

func isValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) != 10 {
		return false
	}

	for index, value := range phoneNumber {
		if index == 0 {
			if value != '0' {
				return false
			}
		} else if index == 1 {
			if value != '6' {
				return false
			}
		} else if !unicode.IsNumber(value) {
			return false
		}
	}

	return true
}

func CreatePoll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "CreatePoll")
}

func GetPoll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	phoneNumber, token, ok := r.BasicAuth()
	if !ok {
		fmt.Fprintf(w, "No authentication provided.")
		return
	}

	user, err := pollyDb.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		fmt.Fprintf(w, "No such user.")
		return
	}

	if user.Token != token {
		fmt.Fprintf(w, "Wrong token.")
		return
	}

	idString := p.ByName("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Fprintf(w, "Bad id.")
	}

	_, err = pollyDb.RetrievePollData(id)
	if err != nil {
		fmt.Fprintf(w, "No poll with that id.")
	}

	fmt.Fprintf(w, "POLL")
}

func GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "GetUser\n")
}

func Vote(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "Vote\n")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
