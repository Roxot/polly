package http

import (
	"fmt"
	"net/http"
	"polly"
	"strconv"
)

func (server *sServer) authenticateRequest(request *http.Request) (
	*polly.PrivateUser, error) {

	idStr, token, ok := request.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("No authentication provided.")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Bad ID.")
	}

	user, err := server.db.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("Unknown user: %s.", idStr)
	}

	if user.Token != token {
		return nil, fmt.Errorf("Bad token.")
	}

	return user, nil
}

func (server *sServer) hasPollAccess(userID int64, pollID int64) bool {
	exists, err := server.db.ExistsParticipant(userID, pollID)
	if err != nil {
		server.logger.Log("hasPollAccess",
			"Somehow existsParticipant returned an error", "hasPollAccess")
		return false
	}

	return exists
}
