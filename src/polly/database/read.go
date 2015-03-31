package database

import "fmt"

import _ "github.com/lib/pq"

func (db *Database) UserByPhoneNumber(phoneNumber string) (PrivateUser, error) {
	var user PrivateUser
	err := db.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName,
			cPhoneNumber), phoneNumber)
	return user, err
}

func (db *Database) UserById(id int) (PrivateUser, error) {
	var user PrivateUser
	err := db.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cId), id)
	return user, err
}

func (db *Database) VerTokenByPhoneNumber(phoneNumber string) (VerToken,
	error) {

	var vt VerToken
	err := db.dbMap.SelectOne(&vt,
		fmt.Sprintf("select * from %s where %s=$1;",
			cVerificationTokensTableName, cPhoneNumber), phoneNumber)
	return vt, err
}

func (db *Database) PollById(id int) (Poll, error) {
	var poll Poll
	err := db.dbMap.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cId), id)
	return poll, err
}

func (db *Database) PollsByUserId(userId int) ([]Poll, error) {
	var polls []Poll
	_, err := db.dbMap.Select(&polls,
		fmt.Sprintf("select id from %s where %s=$1;", cPollTableName,
			cCreatorId), userId)
	return polls, err
}

func (db *Database) QuestionsByPollId(pollId int) ([]Question, error) {
	var questions []Question
	_, err := db.dbMap.Select(&questions,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cPollId), pollId)
	return questions, err
}

func (db *Database) OptionsByPollId(pollId int) ([]Option, error) {
	var options []Option
	_, err := db.dbMap.Select(&options,
		fmt.Sprintf("select * from %s where %s = $1;", cOptionTableName,
			cPollId), pollId)
	return options, err
}

func (db *Database) ParticipantsByPollId(pollId int) (
	[]Participant, error) {

	var participants []Participant
	_, err := db.dbMap.Select(&participants,
		fmt.Sprintf("select * from %s where %s = $1;", cParticipantTableName,
			cPollId), pollId)
	return participants, err
}

func (db *Database) VotesByPollId(pollId int) ([]Vote, error) {
	var votes []Vote
	_, err := db.dbMap.Select(&votes,
		fmt.Sprintf("select * from %s where %s = $1;", cVoteTableName, cPollId),
		pollId)
	return votes, err
}
