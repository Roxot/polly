package http

const (
	cClosedPollEvent = "closed_poll_event"
)

// TODO t -> s
type tPollToClose struct { // TODO should this be here?
	ID    int64
	Title string
}

func (server *sServer) ClosePoll(poll *tPollToClose) error {

	// send a notification to other participants
	err := server.pushClient.NotifyForClosedEvent(&server.db, poll.ID,
		poll.Title)
	if err != nil {
		// TODO neaten up
		server.logger.Log(cClosedPollEvent, "Error notifying: "+err.Error(),
			"::1")
	}

	return nil
}
