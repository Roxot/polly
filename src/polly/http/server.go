package http

import (
	"fmt"
	"net/http"
	"polly/database"
	"polly/log"
	"polly/push"

	"polly/internal/github.com/julienschmidt/httprouter"
)

const (
	cHTTPServerTag  = "HTTPSERVER"
	cAPIVersion     = "v0.1"
	cEndpointFormat = "/api/%s/%s.json"
	// cEndpointWithVarFormat = cEndpointFormat + ":%s"
)

type IServer interface {
	Start(port string) error
	// Stop() TODO ?
}

type sServer struct {
	db         database.Database
	router     httprouter.Router
	logger     log.ILogger
	pushClient push.IPushClient
}

func NewServer(dbConfig *database.DBConfig, clearDB bool) (IServer, error) {
	var err error
	server := sServer{}

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	if clearDB {
		err = db.DropTablesIfExists()
		if err != nil {
			return nil, err
		}
	}

	err = db.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	pushClient, err := push.NewClient()
	if err != nil {
		return nil, err
	}

	server.pushClient = pushClient
	server.logger = log.NewLogger()
	server.db = *db
	server.router = *httprouter.New()

	// start the push notification server's error logging
	err = pushClient.StartErrorLogger(server.logger)
	if err != nil {
		return nil, err
	}

	return &server, nil
}

// is sync
func (server *sServer) Start(port string) error {
	var err error
	err = server.logger.Start()
	if err != nil {
		return err
	}

	server.router.POST(fmt.Sprintf(cEndpointFormat, cAPIVersion, "register"),
		server.Register)
	server.router.PUT(fmt.Sprintf(cEndpointFormat, cAPIVersion, "user"),
		server.UpdateUser)
	server.router.GET(fmt.Sprintf(cEndpointFormat, cAPIVersion, "list_polls"),
		server.ListPolls)
	server.router.POST(fmt.Sprintf(cEndpointFormat, cAPIVersion, "poll"),
		server.PostPoll)
	server.router.GET(fmt.Sprintf(cEndpointFormat, cAPIVersion, "polls"),
		server.GetPollBulk)

	// server.router.GET(fmt.Sprintf(cEndpointFormat
	// 	"/api/v1/user/lookup/:%s", cEmail),
	// 	server.GetUser)
	// server.router.POST("/api/v1/vote", server.Vote)
	// server.router.GET(fmt.Sprintf("/api/v1/poll/:%s", cID), server.GetPoll) TODO deprecated, remove
	server.logger.Log(cHTTPServerTag, "Starting HTTP server", "::1")
	err = http.ListenAndServe(port, &server.router)
	return err
}
