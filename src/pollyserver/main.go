package main

import (
	"fmt"
	"net/http"
	"pollydatabase"

	"github.com/julienschmidt/httprouter"
)

var pollyDb pollydatabase.PollyDatabase

func main() {
	pollyDb, err := pollydatabase.New()
	checkErr(err)

	err = pollyDb.DropTables()
	checkErr(err)

	err = pollyDb.CreateTables()
	checkErr(err)

	newUser := pollydatabase.User{}
	newUser.PhoneNumber = "0612345678"
	newUser.Token = "P0lly"
	newUser.DisplayName = "The Wolf of Wall Street"
	newUser.DeviceType = 0
	newUser.DeviceGUID = ""
	err = pollyDb.AddUser(&newUser)
	checkErr(err)

	// THIS WORKS FINE
	fmt.Println(pollyDb.FindUserByPhoneNumber("0612345668"))

	// THIS DOESNT
	test()

	// router := httprouter.New()
	// router.GET("/poll/:id", GetPoll)
	// router.POST("/poll/:id/vote", Vote)
	// router.POST("/create/poll", CreatePoll)
	// router.POST("/create/user", CreateUser)
	// router.GET("/user/:id/", GetUser)
	//err = http.ListenAndServe(":8081", router)
	//checkErr(err)
}

func test() (pollydatabase.User, error) {
	return pollyDb.FindUserByPhoneNumber("0612345668")
}

func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "CreateUser")
}

func CreatePoll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "CreatePoll")
}

func GetPoll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//fmt.Println(pollyDb.FindUserByPhoneNumber("0612345668"))
	/*
		phoneNumber, _, ok := r.BasicAuth()
		if !ok {
			fmt.Fprintf(w, "No auth provided")
			return
		}

		fmt.Printf("%T, %v", phoneNumber, phoneNumber)

		_, err := pollyDb.FindUserByPhoneNumber("0612345667")
		if err != nil {
			fmt.Fprintf(w, "No such user")
			return
		}*/

	/*
		if user.Token != token {
			fmt.Fprintf(w, "Wrong token")
			return
		}

		fmt.Fprintf(w, "POLL")
	*/
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
