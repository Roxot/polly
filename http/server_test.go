package http

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/roxot/polly"
)

var (
	testUserToken    string
	testUserID       int64
	testUserIDString string
	testEndpoint     string
)

const (
	testConfig = "testing-config.json"
)

func TestMain(m *testing.M) {
	flag.Parse()

	pollyHome, err := polly.GetPollyHome()
	if err != nil {
		log.Fatal(err)
	}

	config, err := ConfigFromFile(filepath.Join(pollyHome, testConfig))
	if err != nil {
		log.Fatal(err)
	}

	server, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	testEndpoint = "http://localhost" + config.Port
	testUserID, testUserToken, err = registerTestUser()
	if err != nil {
		log.Fatal(err)
	}

	// convert the user ID to a string for HTTP requests
	testUserIDString = strconv.FormatInt(testUserID, 10)

	os.Exit(m.Run())
}

func registerTestUser() (int64, string, error) {
	testUser := polly.PrivateUser{
		DisplayName: "Test User",
		DeviceType:  polly.DEVICE_TYPE_ANDROID,
	}

	// marshal the dummy user
	body, err := json.Marshal(&testUser)
	if err != nil {
		return 0, "", err
	}

	// register the dummy user
	reader := bytes.NewReader(body)
	response, err := http.Post(testEndpoint+fmt.Sprintf(cEndpointFormat,
		cAPIVersion, "register"), "application/json", reader)
	if err != nil {
		return 0, "", err
	}

	// unmarshal the response
	var user polly.PrivateUser
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&user); err != nil {
		return 0, "", err
	}

	return user.ID, user.Token, nil
}

func createAuthenticatedRequest(method, api string, body []byte) (*http.Request,
	error) {

	reader := bytes.NewReader(body)
	request, err := http.NewRequest(method, testEndpoint+api, reader)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(testUserIDString, testUserToken)
	return request, nil
}
