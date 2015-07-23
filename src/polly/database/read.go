package database

import (
	"fmt"
	"polly"

	_ "github.com/lib/pq"
)

func (db *Database) GetUserByEmail(email string) (*polly.PrivateUser, error) {
	var user polly.PrivateUser
	err := db.mapping.SelectOne(&user, fmt.Sprintf(
		"select * from %s where %s=$1;", cUserTableName, cEmail), email)
	return &user, err
}

func (db *Database) GetUserByID(id int) (*polly.PrivateUser, error) {
	var user polly.PrivateUser
	err := db.mapping.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cID), id)
	return &user, err
}

func (db *Database) GetPublicUserByEmail(email string) (*polly.PublicUser,
	error) {

	publicUser := polly.PublicUser{}
	privateUser, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	publicUser.ID = privateUser.ID
	publicUser.DisplayName = privateUser.DisplayName
	return &publicUser, nil
}

func (db *Database) GetPublicUserByID(id int) (*polly.PublicUser, error) {
	publicUser := polly.PublicUser{}
	privateUser, err := db.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	publicUser.ID = privateUser.ID
	publicUser.DisplayName = privateUser.DisplayName
	return &publicUser, nil
}

func (db *Database) GetOptionByID(id int) (*polly.Option, error) {
	var option polly.Option
	err := db.mapping.SelectOne(&option,
		fmt.Sprintf("select * from %s where %s=$1;", cOptionTableName, cID), id)
	return &option, err
}

func (db *Database) GetDeviceInfosForPollExcludeCreator(pollID int,
	creatorID int) ([]polly.DeviceInfo, error) {

	var deviceInfos []polly.DeviceInfo
	_, err := db.mapping.Select(&deviceInfos, fmt.Sprintf("select %s.%s, %s.%s"+
		" from %s, %s where %s.%s=%s.%s and %s.%s=$1 and %s.%s!=$2;",
		cUserTableName, cDeviceType, cUserTableName, cDeviceGUID,
		cUserTableName, cParticipantTableName, cUserTableName, cID,
		cParticipantTableName, cUserID, cParticipantTableName, cPollID,
		cUserTableName, cID), pollID, creatorID)

	return deviceInfos, err
}

func (db *Database) GetVerTokenByEmail(email string) (*polly.VerToken, error) {
	var verToken polly.VerToken
	err := db.mapping.SelectOne(&verToken, fmt.Sprintf(
		"select * from %s where %s=$1;", cVerificationTokensTableName,
		cEmail), email)
	return &verToken, err
}

func (db *Database) GetPollByID(id int) (*polly.Poll, error) {
	var poll polly.Poll
	err := db.mapping.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cID), id)
	return &poll, err
}

/*
 * Returns the ordered-by-last-updated list of poll snapshots. Limits the
 * results on the given limit and from the given offset.
 */
func (db *Database) GetPollSnapshotsByUserID(userID, limit, offset int) (
	[]polly.PollSnapshot, error) {

	var snapshots []polly.PollSnapshot
	_, err := db.mapping.Select(&snapshots, fmt.Sprintf(
		"select %s.%s, %s.%s from %s, %s where %s.%s=%s.%s and %s.%s=$1 order"+
			" by %s desc limit %d offset %d;",
		cParticipantTableName, cPollID, cPollTableName, cLastUpdated,
		cParticipantTableName, cPollTableName, cParticipantTableName,
		cPollID, cPollTableName, cID, cParticipantTableName, cUserID,
		cLastUpdated, limit, offset), userID)
	return snapshots, err
}

func (db *Database) GetPollsByUserID(userID int) ([]polly.Poll, error) {
	var polls []polly.Poll
	_, err := db.mapping.Select(&polls,
		fmt.Sprintf("select id from %s where %s=$1;", cPollTableName,
			cCreatorID), userID)
	return polls, err
}

func (db *Database) GetPollIDForOptionID(optionID int) (int, error) {
	var option polly.Option
	err := db.mapping.SelectOne(&option,
		fmt.Sprintf("select %s from %s where %s = $1;", cPollID,
			cOptionTableName, cID), optionID)
	return option.PollID, err
}

func (db *Database) GetPollIDForQuestionID(questionID int) (int, error) {
	var question polly.Question
	err := db.mapping.SelectOne(&question,
		fmt.Sprintf("select %s from %s where %s = $1;", cPollID,
			cQuestionTableName, cID), questionID)
	return question.PollID, err
}

func (db *Database) GetQuestionByPollID(pollID int) (*polly.Question, error) {
	var question polly.Question
	err := db.mapping.SelectOne(&question,
		fmt.Sprintf("select * from %s where %s = $1;", cQuestionTableName,
			cPollID), pollID)
	return &question, err
}

func (db *Database) GetOptionsByPollID(pollID int) ([]polly.Option, error) {
	var options []polly.Option
	_, err := db.mapping.Select(&options,
		fmt.Sprintf("select * from %s where %s = $1;", cOptionTableName,
			cPollID), pollID)
	return options, err
}

func (db *Database) GetParticipantsByPollID(pollID int) (
	[]polly.Participant, error) {

	var participants []polly.Participant
	_, err := db.mapping.Select(&participants,
		fmt.Sprintf("select * from %s where %s = $1;", cParticipantTableName,
			cPollID), pollID)
	return participants, err
}

func (db *Database) GetVotesByPollId(pollID int) ([]polly.Vote, error) {
	var votes []polly.Vote
	_, err := db.mapping.Select(&votes,
		fmt.Sprintf("select * from %s where %s = $1;", cVoteTableName, cPollID),
		pollID)
	return votes, err
}
