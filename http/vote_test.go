package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/roxot/polly"
)

// func BenchmarkVote(b *testing.B) {
// }

func insertDummyPoll() (*polly.PollMessage, error) {
	pollMsg := polly.PollMessage{
		MetaData: polly.Poll{
			ClosingDate: time.Now().Add(time.Hour).UnixNano() / 1000000,
		},
		Question: polly.Question{
			Type:  polly.QUESTION_TYPE_MC,
			Title: "What is the answer?",
		},
		Options: []polly.Option{
			polly.Option{
				Value: "This is the answer",
			},
		},
		Participants: []polly.PublicUser{
			polly.PublicUser{
				ID: testUserID,
			},
		},
	}

	body, err := json.Marshal(&pollMsg)
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf(cEndpointFormat, cAPIVersion, "poll")
	request, err := createAuthenticatedRequest("POST", api, body)
	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&pollMsg); err != nil {
		return nil, err
	}
	return &pollMsg, nil
}
