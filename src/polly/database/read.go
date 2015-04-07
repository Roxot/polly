package database

import (
	"fmt"
	"polly"
)

import _ "github.com/lib/pq"

func (db *Database) UserByPhoneNumber(phoneNumber string) (polly.PrivateUser,
	error) {

	var user polly.PrivateUser
	err := db.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName,
			cPhoneNumber), phoneNumber)
	return user, err
}

func (db *Database) UserById(id int) (polly.PrivateUser, error) {
	var user polly.PrivateUser
	err := db.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cId), id)
	return user, err
}

func (db *Database) PublicUserByPhoneNumber(phoneNumber string) (
	polly.PublicUser, error) {

	pubUser := polly.PublicUser{}
	user, err := db.UserByPhoneNumber(phoneNumber)
	if err != nil {
		return pubUser, err
	}

	pubUser.Id = user.Id
	pubUser.PhoneNumber = user.PhoneNumber
	pubUser.DisplayName = user.DisplayName
	return pubUser, nil
}

func (db *Database) PublicUserById(id int) (polly.PublicUser, error) {
	pubUser := polly.PublicUser{}
	user, err := db.UserById(id)
	if err != nil {
		return pubUser, err
	}

	pubUser.Id = user.Id
	pubUser.PhoneNumber = user.PhoneNumber
	pubUser.DisplayName = user.DisplayName
	return pubUser, nil
}

func (db *Database) VerTokenByPhoneNumber(phoneNumber string) (polly.VerToken,
	error) {

	var vt polly.VerToken
	err := db.dbMap.SelectOne(&vt,
		fmt.Sprintf("select * from %s where %s=$1;",
			cVerificationTokensTableName, cPhoneNumber), phoneNumber)
	return vt, err
}

func (db *Database) PollById(id int) (polly.Poll, error) {
	var poll polly.Poll
	err := db.dbMap.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cId), id)
	return poll, err
}

func (db *Database) PollsByUserId(userId int) ([]polly.Poll, error) {
	var polls []polly.Poll
	_, err := db.dbMap.Select(&polls,
		fmt.Sprintf("select id from %s where %s=$1;", cPollTableName,
			cCreatorId), userId)
	return polls, err
}

func (db *Database) QuestionsByPollId(pollId int) ([]polly.Question, error) {
	var questions []polly.Question
	_, err := db.dbMap.Select(&questions,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cPollId), pollId)
	return questions, err
}

func (db *Database) OptionsByPollId(pollId int) ([]polly.Option, error) {
	var options []polly.Option
	_, err := db.dbMap.Select(&options,
		fmt.Sprintf("select * from %s where %s = $1;", cOptionTableName,
			cPollId), pollId)
	return options, err
}

func (db *Database) ParticipantsByPollId(pollId int) (
	[]polly.Participant, error) {

	var participants []polly.Participant
	_, err := db.dbMap.Select(&participants,
		fmt.Sprintf("select * from %s where %s = $1;", cParticipantTableName,
			cPollId), pollId)
	return participants, err
}

func (db *Database) VotesByPollId(pollId int) ([]polly.Vote, error) {
	var votes []polly.Vote
	_, err := db.dbMap.Select(&votes,
		fmt.Sprintf("select * from %s where %s = $1;", cVoteTableName, cPollId),
		pollId)
	return votes, err
}
