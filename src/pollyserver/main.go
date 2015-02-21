package main

import (
	"log"
	"polly/database"
	"polly/httpserver"
)

const (
	cPsqlUser = "polly"
	cPsqlPass = "w01V3s"
	cDbName   = "pollydb"
	cPort     = ":8080"
	cClearDb  = false
)

func main() {

	dbConfig := database.DbConfig{}
	dbConfig.PsqlUser = cPsqlUser
	dbConfig.PsqlUserPass = cPsqlPass
	dbConfig.DbName = cDbName

	log.Println("Opening database...")
	srv, err := httpserver.New(dbConfig, cClearDb)
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
