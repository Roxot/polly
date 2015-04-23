package httpserver

import (
	"encoding/json"
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
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type VoteResponseMessage struct {
	Option *polly.Option `json:"option,omitempty"`
	Vote   *polly.Vote   `json:"vote"`
}

func (srv *HTTPServer) Vote(writer http.ResponseWriter, req *http.Request,
	_ httprouter.Params) {

	// authenticate the user
	usr, err := srv.authenticateRequest(req)
	if err != nil {
		srv.handleAuthError(cVoteTag, err, writer, req)
		return
	}

	// decode the vote message
	var voteMsg VoteMessage
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&voteMsg)
	if err != nil {
		srv.handleBadRequest(cVoteTag, cBadJSONErr, err, writer, req)
		return
	}

	// retrieve the poll id belonging to the option or question id
	var pollId int
	switch voteMsg.Type {
	case VOTE_TYPE_NEW:
		pollId, err = srv.db.PollIdForQuestionId(voteMsg.Id)
		if err != nil {
			srv.handleBadRequest(cVoteTag, cNoQuestionErr, err, writer, req)
			return
		} else if len(voteMsg.Value) == 0 {
			srv.handleErr(cVoteTag, cBadVoteMsgErr, cBadVoteMsgErr, 400,
				writer, req)
			return
		}
	case VOTE_TYPE_UPVOTE:
		pollId, err = srv.db.PollIdForOptionId(voteMsg.Id)
		if err != nil {
			srv.handleBadRequest(cVoteTag, cNoOptionErr, err, writer, req)
			return
		}
	default:
		srv.handleErr(cVoteTag, cBadVoteTypeErr, cBadVoteTypeErr, 400,
			writer, req)
		return
	}

	// make sure the user is allowed to vote
	if !srv.hasPollAccess(usr.Id, pollId) {
		srv.handleIllegalOperation(cVoteTag, cAccessRightsErr, writer, req)
		return
	}

	// start a transaction
	transaction, err := srv.db.Begin()
	if err != nil {
		transaction.Rollback()
		srv.handleDatabaseError(cVoteTag, err, writer, req)
		return
	}

	// remove all existing votes of the user
	err = database.DelVotesForUserTx(usr.Id, pollId, transaction)
	if err != nil {
		transaction.Rollback()
		srv.handleDatabaseError(cVoteTag, err, writer, req)
		return
	}

	// if necessary, create a new option
	var optionId int
	var option polly.Option
	if voteMsg.Type == VOTE_TYPE_UPVOTE {
		optionId = voteMsg.Id
	} else {

		// we have a vote message with type NEW, so we create a new option
		questionId := voteMsg.Id
		option.PollId = pollId
		option.QuestionId = questionId
		option.Value = voteMsg.Value
		err = srv.db.AddOptionTx(&option, transaction)
		if err != nil {
			transaction.Rollback()
			srv.handleDatabaseError(cVoteTag, err, writer, req)
			return
		}

		optionId = option.Id
	}

	// insert the vote into the database
	vote := polly.Vote{}
	vote.CreationDate = time.Now().Unix()
	vote.OptionId = optionId
	vote.PollId = pollId
	vote.UserId = usr.Id
	err = database.AddVoteTx(&vote, transaction)
	if err != nil {
		transaction.Rollback()
		srv.handleDatabaseError(cVoteTag, err, writer, req)
		return
	}

	// update the poll last updated
	err = database.UpdatePollLastUpdatedTx(pollId, time.Now().Unix(),
		transaction)
	if err != nil {
		transaction.Rollback()
		srv.handleDatabaseError(cVoteTag, err, writer, req)
		return
	}

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		srv.handleDatabaseError(cVoteTag, err, writer, req)
		return
	}

	// construct the response message
	response := VoteResponseMessage{}
	response.Vote = &vote
	if voteMsg.Type == VOTE_TYPE_NEW {
		response.Option = &option
	}

	// send the response message
	responseBody, err := json.MarshalIndent(response, "", "\t")
	_, err = writer.Write(responseBody)
	if err != nil {
		srv.handleWritingError(cVoteTag, err, writer, req)
		return
	}

}
