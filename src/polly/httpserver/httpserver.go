package httpserver

import (
	"net/http"
	"polly/database"

	"github.com/julienschmidt/httprouter"
)

type HTTPServer struct {
	db     database.Database
	router *httprouter.Router
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

	srv.db = db
	srv.router = httprouter.New()

	return srv, nil
}

func (srv *HTTPServer) Database() *database.Database {
	return &srv.db
}

// sync
func (srv *HTTPServer) Start(port string) error {
	srv.router.POST("/register", srv.Register)
	srv.router.POST("/register/verify", srv.VerifyRegister)
	err := http.ListenAndServe(port, srv.router)
	return err
}
