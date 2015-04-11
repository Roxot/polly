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

	dbCfg := database.DbConfig{}
	dbCfg.PsqlUser = cPsqlUser
	dbCfg.PsqlUserPass = cPsqlPass
	dbCfg.DbName = cDbName

	log.Println("Opening database...")
	srv, err := httpserver.New(&dbCfg, cClearDb)
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
