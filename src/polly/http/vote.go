package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"polly"
	"polly/database"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	VOTE_TYPE_NEW    = 0
	VOTE_TYPE_UPVOTE = 1

	cVoteTag = "POST/VOTE"
)

type VoteMessage struct {
	Type  int    `json:"type"`
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type VoteResponseMessage struct {
	Option polly.Option `json:"option,omitempty"`
	Vote   polly.Vote   `json:"vote"`
}

func (server *sServer) Vote(writer http.ResponseWriter, request *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	user, err := server.authenticateRequest(request)
	if err != nil {
		server.handleAuthError(cVoteTag, err, writer, request)
		return
	}

	// decode the vote message
	var voteMsg VoteMessage
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&voteMsg)
	if err != nil {
		server.handleBadRequest(cVoteTag, cBadJSONErr, err, writer, request)
		return
	}

	// retrieve the poll id belonging to the option or question id
	var pollID int
	switch voteMsg.Type {
	case VOTE_TYPE_NEW:
		pollID, err = server.db.GetPollIDForQuestionID(voteMsg.ID)
		if err != nil {
			server.handleBadRequest(cVoteTag, cNoQuestionErr, err, writer,
				request)
			return
		} else if len(voteMsg.Value) == 0 {
			server.handleErr(cVoteTag, cBadVoteMsgErr, cBadVoteMsgErr, 400,
				writer, request)
			return
		}
	case VOTE_TYPE_UPVOTE:
		pollID, err = server.db.GetPollIDForOptionID(voteMsg.ID)
		if err != nil {
			server.handleBadRequest(cVoteTag, cNoOptionErr, err, writer,
				request)
			return
		}
	default:
		server.handleErr(cVoteTag, cBadVoteTypeErr, cBadVoteTypeErr, 400,
			writer, request)
		return
	}

	// make sure the user is allowed to vote
	if !server.hasPollAccess(user.ID, pollID) {
		server.handleIllegalOperation(cVoteTag, cAccessRightsErr, writer,
			request)
		return
	}

	// start a transaction
	transaction, err := server.db.Begin()
	if err != nil {
		transaction.Rollback()
		server.handleDatabaseError(cVoteTag, err, writer, request)
		return
	}

	// remove all existing votes of the user
	err = database.DeleteVotesForUserTX(user.ID, pollID, transaction)
	if err != nil {
		transaction.Rollback()
		server.handleDatabaseError(cVoteTag, err, writer, request)
		return
	}

	// if necessary, create a new option
	var optionID int
	var option polly.Option
	if voteMsg.Type == VOTE_TYPE_UPVOTE {
		optionID = voteMsg.ID
	} else {

		// we have a vote message with type NEW, so we create a new option
		questionID := voteMsg.ID
		option.PollID = pollID
		option.QuestionID = questionID
		option.Value = voteMsg.Value
		err = database.AddOptionTX(&option, transaction)
		if err != nil {
			transaction.Rollback()
			server.handleDatabaseError(cVoteTag, err, writer, request)
			return
		}

		optionID = option.ID
	}

	// insert the vote into the database
	vote := polly.Vote{}
	vote.CreationDate = time.Now().Unix()
	vote.OptionID = optionID
	vote.PollID = pollID
	vote.UserID = user.ID
	err = database.AddVoteTX(&vote, transaction)
	if err != nil {
		transaction.Rollback()
		server.handleDatabaseError(cVoteTag, err, writer, request)
		return
	}

	// update the poll last updated
	err = database.UpdatePollLastUpdatedTX(pollID, time.Now().Unix(),
		transaction)
	if err != nil {
		transaction.Rollback()
		server.handleDatabaseError(cVoteTag, err, writer, request)
		return
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		server.handleDatabaseError(cVoteTag, err, writer, request)
		return
	}

	// push a notification to all participants of the poll TODO other type of vote
	err = server.pushClient.NotifyForUpvote(&server.db, user, optionID)
	if err != nil {
		host, _, _ := net.SplitHostPort(request.RemoteAddr)
		server.logger.Log(cVoteTag, fmt.Sprintf(cDatabaseErrLog, err), host)
	}

	// construct the response message
	response := VoteResponseMessage{}
	response.Vote = vote
	if voteMsg.Type == VOTE_TYPE_NEW {
		response.Option = option
	}

	// send the response message
	responseBody, err := json.MarshalIndent(response, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		server.handleWritingError(cVoteTag, err, writer, request)
		return
	}

}
