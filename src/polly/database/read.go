package database

import (
	"fmt"
	"polly"

	_ "polly/internal/github.com/lib/pq"
	"polly/internal/gopkg.in/gorp.v1"
)

func (db *Database) GetUserByID(id int64) (*polly.PrivateUser, error) {
	var user polly.PrivateUser
	err := db.mapping.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cID), id)
	return &user, err
}

func (db *Database) GetPublicUserByID(id int64) (*polly.PublicUser, error) {
	publicUser := polly.PublicUser{}
	privateUser, err := db.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	publicUser.ID = privateUser.ID
	publicUser.DisplayName = privateUser.DisplayName
	return &publicUser, nil
}

func (db *Database) GetOptionByID(id int64) (*polly.Option, error) {
	var option polly.Option
	err := db.mapping.SelectOne(&option,
		fmt.Sprintf("select * from %s where %s=$1;", cOptionTableName, cID), id)
	return &option, err
}

func (db *Database) GetDeviceInfosForPollExcludeCreator(pollID,
	creatorID int64) ([]polly.DeviceInfo, error) {

	var deviceInfos []polly.DeviceInfo
	_, err := db.mapping.Select(&deviceInfos, fmt.Sprintf("select %s.%s, %s.%s"+
		" from %s, %s where %s.%s=%s.%s and %s.%s=$1 and %s.%s!=$2;",
		cUserTableName, cDeviceType, cUserTableName, cDeviceGUID,
		cUserTableName, cParticipantTableName, cUserTableName, cID,
		cParticipantTableName, cUserID, cParticipantTableName, cPollID,
		cUserTableName, cID), pollID, creatorID)

	return deviceInfos, err
}

func (db *Database) GetDeviceInfosForPoll(pollID int64) ([]polly.DeviceInfo, 
	error) {

	var deviceInfos []polly.DeviceInfo
	_, err := db.mapping.Select(&deviceInfos, fmt.Sprintf("select %s.%s, %s.%s"+
		" from %s, %s where %s.%s=%s.%s and %s.%s=$1;",
		cUserTableName, cDeviceType, cUserTableName, cDeviceGUID,
		cUserTableName, cParticipantTableName, cUserTableName, cID,
		cParticipantTableName, cUserID, cParticipantTableName, cPollID), 
		pollID)

	return deviceInfos, err
}


func (db *Database) GetPollByID(id int64) (*polly.Poll, error) {
	var poll polly.Poll
	err := db.mapping.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cID), id)
	return &poll, err
}

/*
 * Returns the ordered-by-last-updated list of poll snapshots. Limits the
 * results on the given limit and from the given offset.
 */
func (db *Database) GetPollSnapshotsByUserID(userID int64, limit, offset int) (
	[]polly.PollSnapshot, error) {

	var snapshots []polly.PollSnapshot
	_, err := db.mapping.Select(&snapshots, fmt.Sprintf(
		"select %s.%s, %s.%s, %s.%s, %s.%s from %s, %s where %s.%s=%s.%s and "+
			"%s.%s=$1 order by %s desc limit %d offset %d;",
		cPollTableName, cID, cPollTableName, cLastUpdated,
		cPollTableName, cSequenceNumber, cPollTableName, cClosingDate, 
		cParticipantTableName, cPollTableName, cParticipantTableName, cPollID,
		cPollTableName, cID, cParticipantTableName, cUserID, cLastUpdated,
		limit, offset), userID)
	return snapshots, err
}

func (db *Database) GetPollsByUserID(userID int64) ([]polly.Poll, error) {
	var polls []polly.Poll
	_, err := db.mapping.Select(&polls,
		fmt.Sprintf("select id from %s where %s=$1;", cPollTableName,
			cCreatorID), userID)
	return polls, err
}

func (db *Database) GetPollIDForOptionID(optionID int64) (int64, error) {
	var option polly.Option
	err := db.mapping.SelectOne(&option,
		fmt.Sprintf("select %s from %s where %s=$1;", cPollID,
			cOptionTableName, cID), optionID)
	return option.PollID, err
}

func (db *Database) GetPollIDForQuestionID(questionID int64) (int64, error) {
	var question polly.Question
	err := db.mapping.SelectOne(&question,
		fmt.Sprintf("select %s from %s where %s=$1;", cPollID,
			cQuestionTableName, cID), questionID)
	return question.PollID, err
}

func (db *Database) GetPollIDForVoteID(voteID int64) (int64, error) {
	var vote polly.Vote
	err := db.mapping.SelectOne(&vote, // TODO selectInt?
		fmt.Sprintf("select %s from %s where %s=$1;", cPollID,
			cVoteTableName, cID), voteID)
	return vote.PollID, err
}

func (db *Database) GetQuestionByPollID(pollID int64) (*polly.Question, error) {
	var question polly.Question
	err := db.mapping.SelectOne(&question,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cPollID), pollID)
	return &question, err
}

func (db *Database) GetQuestionByID(questionID int64) (*polly.Question, error) {
	var question polly.Question
	err := db.mapping.SelectOne(&question,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cID), questionID)
	return &question, err
}

func (db *Database) GetOptionsByPollID(pollID int64) ([]polly.Option, error) {
	var options []polly.Option
	_, err := db.mapping.Select(&options,
		fmt.Sprintf("select * from %s where %s = $1;", cOptionTableName,
			cPollID), pollID)
	return options, err
}

func (db *Database) GetParticipantsByPollID(pollID int64) (
	[]polly.Participant, error) {

	var participants []polly.Participant
	_, err := db.mapping.Select(&participants,
		fmt.Sprintf("select * from %s where %s = $1;", cParticipantTableName,
			cPollID), pollID)
	return participants, err
}

func (db *Database) GetVotesByPollID(pollID int64) ([]polly.Vote, error) {
	var votes []polly.Vote
	_, err := db.mapping.Select(&votes,
		fmt.Sprintf("select * from %s where %s = $1;", cVoteTableName, cPollID),
		pollID)
	return votes, err
}

func (db *Database) GetSequenceNumber(pollID int64) (int, error) {
	number, err := db.mapping.SelectInt(fmt.Sprintf(
		"select %s from %s where %s=$1;", cSequenceNumber, cPollTableName, cID),
		pollID)
	return int(number), err
}

func (db *Database) GetClosingDate(pollID int64) (int64, error) {
	number, err := db.mapping.SelectInt(fmt.Sprintf(
		"select %s from %s where %s=$1;", cClosingDate, cPollTableName, cID),
		pollID)
	return number, err
}

func GetSequenceNumberTX(pollID int64, tx *gorp.Transaction) (int, error) {
	number, err := tx.SelectInt(fmt.Sprintf(
		"select %s from %s where %s=$1;", cSequenceNumber, cPollTableName, cID),
		pollID)
	return int(number), err
}

func GetPollSnapshotTX(pollID int64, tx *gorp.Transaction) (*polly.PollSnapshot,
	error) {

	var snapshot polly.PollSnapshot
	err := tx.SelectOne(&snapshot, fmt.Sprintf(
		"select %s, %s, %s, %s from %s where %s=$1;", cID, cLastUpdated,
		cSequenceNumber, cClosingDate, cPollTableName, cID), pollID)
	return &snapshot, err
}
