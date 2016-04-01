package http

import (
	"fmt"
	"net/http"

	"github.com/albrow/jobs"
	"github.com/roxot/polly/database"
	"github.com/roxot/polly/log"
	"github.com/roxot/polly/push"

	"github.com/julienschmidt/httprouter"
)

const (
	cHTTPServerTag   = "HTTPSERVER"
	cAPIVersion      = "v0.1"
	cEndpointFormat  = "/%s/%s.json"
	cClosedPollsJobs = "CLOSED_POLLS"
	// cEndpointWithVarFormat = cEndpointFormat + ":%s"
)

type IServer interface {
	Start() error
	// Stop() TODO ?
}

type sServer struct {
	db          database.Database
	router      httprouter.Router
	logger      log.ILogger
	pushClient  push.IPushClient
	cpScheduler jobs.Type
	config      Config
}

func NewServer(config *Config) (IServer, error) {
	var err error
	server := sServer{}

	db, err := database.NewDatabase(&config.DBConfig)
	if err != nil {
		return nil, err
	}

	if config.TruncateDB {
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
	server.config = *config

	// start the push notification server's error logging
	err = pushClient.StartErrorLogger(server.logger)
	if err != nil {
		return nil, err
	}

	// register the closed poll scheduler
	cpScheduler, err := jobs.RegisterType(cClosedPollsJobs,
		config.ClosedPollPushRetries, server.ClosePoll)
	if err != nil {
		return nil, err
	}

	// create a job pool for the closed poll scheduler TODO pass configuration s
	pool, err := jobs.NewPool(nil)
	if err != nil {
		return nil, err
	}

	// start the job pool for the closed poll scheduler
	err = pool.Start()
	if err != nil {
		return nil, err
	}

	server.cpScheduler = *cpScheduler

	return &server, nil
}

// is sync
func (server *sServer) Start() error {
	var err error
	err = server.logger.Start()
	if err != nil {
		return err
	}

	// TODO endpoint formatting to function
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
	server.router.POST(fmt.Sprintf(cEndpointFormat, cAPIVersion, "vote"),
		server.Vote)
	server.router.GET(fmt.Sprintf(cEndpointFormat, cAPIVersion, "users"),
		server.GetUserBulk)
	server.router.DELETE(fmt.Sprintf(cEndpointFormat, cAPIVersion, "vote"),
		server.UndoVote)
	server.router.DELETE(fmt.Sprintf(cEndpointFormat, cAPIVersion, "poll"),
		server.LeavePoll)
	server.router.POST(fmt.Sprintf(cEndpointFormat, cAPIVersion, "adduser"),
		server.AddUser)
	server.logger.Log(cHTTPServerTag, "Starting HTTP server", "::1")
	err = http.ListenAndServe(server.config.Port, &server.router)
	return err
}
