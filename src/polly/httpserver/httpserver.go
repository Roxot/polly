package httpserver

import (
	"fmt"
	"net/http"
	"polly/database"
	"polly/logger"

	"github.com/julienschmidt/httprouter"
)

type HTTPServer struct {
	db     database.Database
	router *httprouter.Router
	logger logger.Logger
}

func New(dbConfig database.DbConfig, clearDb bool) (HTTPServer, error) {
	var err error
	srv := HTTPServer{}

	db, err := database.New(dbConfig)
	if err != nil {
		return srv, err
	}

	if clearDb {
		err = db.DropTablesIfExists()
		if err != nil {
			return srv, err
		}
	}

	err = db.CreateTablesIfNotExists()
	if err != nil {
		return srv, err
	}

	srv.logger = logger.New()
	srv.db = db
	srv.router = httprouter.New()

	return srv, nil
}

func (srv *HTTPServer) Database() *database.Database {
	return &srv.db
}

// sync
func (srv *HTTPServer) Start(port string) error {
	var err error
	err = srv.logger.Start()
	if err != nil {
		return err
	}

	srv.router.POST("/register", srv.Register)
	srv.router.POST("/register/verify", srv.VerifyRegister)
	srv.router.POST("/poll", srv.PostPoll)
	srv.router.GET("/user/polls", srv.ListUserPolls)
	srv.router.GET(fmt.Sprintf("/poll/:%s", cId), srv.GetPoll)
	srv.router.GET("/poll", srv.GetPollBulk)
	err = http.ListenAndServe(port, srv.router)
	return err
}
