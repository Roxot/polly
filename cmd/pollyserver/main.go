package main

import (
	"flag"
	"log"

	"github.com/roxot/polly/database"
	"github.com/roxot/polly/http"
)

const (
	cPsqlUser     = "polly"
	cPsqlPassword = "w01V3s"
	cDBName       = "pollydb"
	cPort         = ":8080"

	cFlagClearDB = "cleardb"
)

func main() {

	var clearDB bool
	flag.BoolVar(&clearDB, cFlagClearDB, false,
		"Set to true to reset the database.")
	flag.Parse()

	dbConfig := database.DBConfig{}
	dbConfig.User = cPsqlUser
	dbConfig.Password = cPsqlPassword
	dbConfig.Name = cDBName

	log.Println("Opening database...")
	srv, err := http.NewServer(&dbConfig, clearDB)
	if err != nil {
		panic(err)
	} else {
		log.Println("Database opened successfully.")
	}

	log.Printf("Starting HTTP server on port %s...\n", cPort)
	err = srv.Start(cPort)
	if err != nil {
		panic(err)
	}

}
