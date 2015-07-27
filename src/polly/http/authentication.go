package http

import (
	"net/http"
	"polly"
)

func (server *sServer) authenticateRequest(request *http.Request) (
	*polly.PrivateUser, error) {

	// email, token, ok := request.BasicAuth()
	// if !ok {
	// 	return nil, fmt.Errorf("No authentication provided.")
	// }

	// user, err := server.db.GetUserByEmail(email)
	// if err != nil {
	// 	return nil, fmt.Errorf("Unknown user: %s.", email)
	// }

	// if user.Token != token {
	// 	return nil, fmt.Errorf("Bad token.")
	// }

	// return user, nil
	return nil, nil // TODO fix
}

func (server *sServer) hasPollAccess(userID int64, pollID int) bool {
	exists, err := server.db.ExistsParticipant(userID, pollID)
	if err != nil {
		server.logger.Log("hasPollAccess",
			"Somehow existsParticipant returned an error", "hasPollAccess")
		return false
	}

	return exists
}
