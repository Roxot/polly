package http

import (
	"net/http"
	"polly"
	"strconv"
)

func (server *sServer) authenticateRequest(request *http.Request) (
	*polly.PrivateUser, int) {

	idStr, token, ok := request.BasicAuth()
	if !ok {
		return nil, ERR_AUT_NO_AUTH
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, ERR_BAD_ID
	}

	user, err := server.db.GetUserByID(id)
	if err != nil {
		return nil, ERR_AUT_NO_USER
	}

	if user.Token != token {
		return nil, ERR_AUT_BAD_TOKEN
	}

	return user, NO_ERR
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
