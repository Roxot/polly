package database

import "fmt"

import _ "github.com/lib/pq"

func (pollyDb Database) FindUserByPhoneNumber(phoneNumber string) (User,
	error) {

	var user User
	err := pollyDb.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName,
			cPhoneNumber), phoneNumber)
	return user, err
}

func (pollyDb Database) FindVerificationTokenByPhoneNumber(
	phoneNumber string) (VerificationToken, error) {

	var vt VerificationToken
	err := pollyDb.dbMap.SelectOne(&vt,
		fmt.Sprintf("select * from %s where %s=$1;",
			cVerificationTokensTableName, cPhoneNumber), phoneNumber)
	return vt, err
}

func (pollyDb Database) FindUserById(id int) (User, error) {
	var user User
	err := pollyDb.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cId), id)
	return user, err
}

func (pollydb Database) FindPollById(id int) (Poll, error) {
	var poll Poll
	err := pollydb.dbMap.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cId), id)
	return poll, err
}

func (pollyDb Database) FindQuestionsByPollId(pollId int) ([]Question,
	error) {

	var questions []Question
	_, err := pollyDb.dbMap.Select(&questions,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cPollId), pollId)
	return questions, err
}

func (pollyDb Database) FindOptionsByPollId(pollId int) ([]Option, error) {
	var options []Option
	_, err := pollyDb.dbMap.Select(&options,
		fmt.Sprintf("select * from %s where %s = $1;", cOptionTableName,
			cPollId), pollId)
	return options, err
}

func (pollyDb Database) FindVotesByPollId(pollId int) ([]Vote, error) {
	var votes []Vote
	_, err := pollyDb.dbMap.Select(&votes,
		fmt.Sprintf("select * from %s where %s = $1;", cVoteTableName, cPollId),
		pollId)
	return votes, err
}

func (pollyDb Database) FindParticipantsByPollId(pollId int) (
	[]Participant, error) {

	var participants []Participant
	_, err := pollyDb.dbMap.Select(&participants,
		fmt.Sprintf("select * from %s where %s = $1;", cParticipantTableName,
			cPollId), pollId)
	return participants, err
}
