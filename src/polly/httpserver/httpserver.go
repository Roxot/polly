package httpserver

import (
	"fmt"
	"net/http"
	"polly/database"
	"polly/httpserver/pushserver"
	"polly/logger"

	"github.com/julienschmidt/httprouter"
)

const cHTTPServerTag = "HTTPSERVER"

type HTTPServer struct {
	db      *database.Database
	router  *httprouter.Router
	logger  *logger.Logger
	pushSrv *pushserver.PushServer
}

func New(dbCfg *database.DbConfig, clearDb bool) (*HTTPServer, error) {
	var err error
	srv := HTTPServer{}

	db, err := database.New(dbCfg)
	if err != nil {
		return nil, err
	}

	if clearDb {
		err = db.DropTablesIfExists()
		if err != nil {
			return nil, err
		}
	}

	err = db.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	pushSrv, err := pushserver.New()
	if err != nil {
		return nil, err
	}

	srv.pushSrv = pushSrv
	srv.logger = logger.New()
	srv.db = db
	srv.router = httprouter.New()

	// start the push notification server's error logging
	err = pushSrv.StartErrorLogger(srv.logger)
	if err != nil {
		return nil, err
	}

	return &srv, nil
}

// sync
func (srv *HTTPServer) Start(port string) error {
	var err error
	err = srv.logger.Start()
	if err != nil {
		return err
	}

	srv.router.POST("/api/v1/register", srv.Register)
	srv.router.POST("/api/v1/register/verify", srv.VerifyRegister)
	srv.router.POST("/api/v1/poll", srv.PostPoll)
	srv.router.POST("/api/v1/vote", srv.Vote)
	srv.router.GET("/api/v1/user/polls", srv.ListUserPolls)
	srv.router.GET(fmt.Sprintf("/api/v1/poll/:%s", cId), srv.GetPoll)
	srv.router.GET("/api/v1/poll", srv.GetPollBulk)
	srv.router.GET(fmt.Sprintf("/api/v1/user/lookup/:%s", cEmail), srv.GetUser)
	srv.logger.Log(cHTTPServerTag, "Starting HTTP server", "::1")
	err = http.ListenAndServe(port, srv.router)
	return err
}
