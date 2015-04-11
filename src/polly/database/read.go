package database

import (
	"fmt"
	"polly"

	_ "github.com/lib/pq"
)

func (db *Database) UserByPhoneNumber(phoneNo string) (*polly.PrivateUser,
	error) {

	var usr polly.PrivateUser
	err := db.dbMap.SelectOne(&usr,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName,
			cPhoneNumber), phoneNo)
	return &usr, err
}

func (db *Database) UserById(id int) (*polly.PrivateUser, error) {
	var usr polly.PrivateUser
	err := db.dbMap.SelectOne(&usr,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cId), id)
	return &usr, err
}

func (db *Database) PublicUserByPhoneNumber(phoneNo string) (
	*polly.PublicUser, error) {

	pubUsr := polly.PublicUser{}
	privUsr, err := db.UserByPhoneNumber(phoneNo)
	if err != nil {
		return nil, err
	}

	pubUsr.Id = privUsr.Id
	pubUsr.PhoneNumber = privUsr.PhoneNumber
	pubUsr.DisplayName = privUsr.DisplayName
	return &pubUsr, nil
}

func (db *Database) PublicUserById(id int) (*polly.PublicUser, error) {
	pubUsr := polly.PublicUser{}
	privUsr, err := db.UserById(id)
	if err != nil {
		return nil, err
	}

	pubUsr.Id = privUsr.Id
	pubUsr.PhoneNumber = privUsr.PhoneNumber
	pubUsr.DisplayName = privUsr.DisplayName
	return &pubUsr, nil
}

func (db *Database) VerTokenByPhoneNumber(phoneNo string) (*polly.VerToken,
	error) {

	var verTkn polly.VerToken
	err := db.dbMap.SelectOne(&verTkn,
		fmt.Sprintf("select * from %s where %s=$1;",
			cVerificationTokensTableName, cPhoneNumber), phoneNo)
	return &verTkn, err
}

func (db *Database) PollById(id int) (*polly.Poll, error) {
	var poll polly.Poll
	err := db.dbMap.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cId), id)
	return &poll, err
}

/*
 * Returns the ordered-by-last-updated list of poll objects with only the id
 * and last updated fields filled in. Limits the results on the given limit
 * and from the given offset.
 */
func (db *Database) PollSnapshotsByUserId(userId, limit, offset int) (
	[]polly.PollSnapshot, error) {

	var snapshots []polly.PollSnapshot
	_, err := db.dbMap.Select(&snapshots, fmt.Sprintf(
		"select %s.%s, %s.%s from %s, %s where %s.%s=%s.%s order by %s desc limit %d offset %d;",
		cParticipantTableName, cPollId, cPollTableName, cLastUpdated,
		cParticipantTableName, cPollTableName, cParticipantTableName,
		cPollId, cPollTableName, cId, cLastUpdated, limit, offset))
	return snapshots, err
}

func (db *Database) PollsByUserId(userId int) ([]polly.Poll, error) {
	var polls []polly.Poll
	_, err := db.dbMap.Select(&polls,
		fmt.Sprintf("select id from %s where %s=$1;", cPollTableName,
			cCreatorId), userId)
	return polls, err
}

func (db *Database) QuestionByPollId(pollId int) (*polly.Question, error) {
	var question polly.Question
	err := db.dbMap.SelectOne(&question,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cPollId), pollId)
	return &question, err
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

	var partics []polly.Participant
	_, err := db.dbMap.Select(&partics,
		fmt.Sprintf("select * from %s where %s = $1;", cParticipantTableName,
			cPollId), pollId)
	return partics, err
}

func (db *Database) VotesByPollId(pollId int) ([]polly.Vote, error) {
	var votes []polly.Vote
	_, err := db.dbMap.Select(&votes,
		fmt.Sprintf("select * from %s where %s = $1;", cVoteTableName, cPollId),
		pollId)
	return votes, err
}
