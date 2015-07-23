package main

import (
	"flag"
	"log"
	"polly/database"
	"polly/http"
)

const (
	cPsqlUser     = "polly"
	cPsqlPassword = "w01V3s"
	cDbName       = "pollydb"
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
	dbConfig.Name = cDbName

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
