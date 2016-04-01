package main

import (
	"fmt"
	"log"
	"os"

	"github.com/roxot/polly/http"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	config, err := http.ConfigFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	srv, err := http.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting HTTP server on port %s...\n", config.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

}

func printUsage() {
	fmt.Println("Usage: pollyserver <config>")
}
