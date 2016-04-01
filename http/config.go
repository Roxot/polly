package http

import (
	"encoding/json"
	"os"

	"github.com/roxot/polly/database"
)

type Config struct {
	DBConfig                   database.Config
	TruncateDB                 bool
	Port                       string
	ClosedPollPushRetries      uint
	VerifyRegisterWithFacebook bool
}

func ConfigFromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
