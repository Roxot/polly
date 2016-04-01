package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/roxot/polly"
)

func BenchmarkPostPoll(b *testing.B) {
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
		b.Fatal(err)
	}

	api := fmt.Sprintf(cEndpointFormat, cAPIVersion, "poll")
	client := new(http.Client)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request, err := createAuthenticatedRequest("POST", api, body)
		if err != nil {
			b.Fatal(err)
		}

		response, err := client.Do(request)
		if err != nil || response.StatusCode != http.StatusOK {
			b.Fatal(err)
		}
	}
}
