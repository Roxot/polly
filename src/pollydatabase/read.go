package pollydatabase

import "fmt"

import _ "github.com/lib/pq"

type PollData struct {
	Poll         Poll
	Questions    []Question
	Options      []Option
	Votes        []Vote
	Creator      User
	Participants []User
}

func (pollyDb PollyDatabase) FindUserByPhoneNumber(phoneNumber string) (User,
	error) {

	var user User
	err := pollyDb.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName,
			cPhoneNumber), phoneNumber)
	return user, err
}

func (pollyDb PollyDatabase) FindVerificationTokenByPhoneNumber(
	phoneNumber string) (VerificationToken, error) {

	var vt VerificationToken
	err := pollyDb.dbMap.SelectOne(&vt,
		fmt.Sprintf("select * from %s where %s=$1;",
			cVerificationTokensTableName, cPhoneNumber), phoneNumber)
	return vt, err
}

func (pollyDb PollyDatabase) FindUserById(id int) (User, error) {
	var user User
	err := pollyDb.dbMap.SelectOne(&user,
		fmt.Sprintf("select * from %s where %s=$1;", cUserTableName, cId), id)
	return user, err
}

func (pollydb PollyDatabase) FindPollById(id int) (Poll, error) {
	var poll Poll
	err := pollydb.dbMap.SelectOne(&poll,
		fmt.Sprintf("select * from %s where %s=$1;", cPollTableName, cId), id)
	return poll, err
}

func (pollyDb PollyDatabase) FindQuestionsByPollId(pollId int) ([]Question,
	error) {

	var questions []Question
	_, err := pollyDb.dbMap.Select(&questions,
		fmt.Sprintf("select * for %s where %s = $1;", cQuestionTableName,
			cPollId), pollId)
	return questions, err
}

func (pollyDb PollyDatabase) FindOptionsByPollId(pollId int) ([]Option, error) {
	var options []Option
	_, err := pollyDb.dbMap.Select(&options,
		fmt.Sprintf("select * for %s where %s = $1;", cOptionTableName,
			cPollId), pollId)
	return options, err
}

func (pollyDb PollyDatabase) FindVotesByPollId(pollId int) ([]Vote, error) {
	var votes []Vote
	_, err := pollyDb.dbMap.Select(&votes,
		fmt.Sprintf("select * for %s where %s = $1;", cVoteTableName, cPollId),
		pollId)
	return votes, err
}

func (pollyDb PollyDatabase) FindParticipantsByPollId(pollId int) (
	[]Participant, error) {

	var participants []Participant
	_, err := pollyDb.dbMap.Select(&participants,
		fmt.Sprintf("select * for %s where %s = $1;", cParticipantTableName,
			cPollId), pollId)
	return participants, err
}

func (pollyDb PollyDatabase) RetrievePollData(pollId int) (PollData, error) {
	pollData := PollData{}

	poll, err := pollyDb.FindPollById(pollId)
	pollData.Poll = poll
	if err != nil {
		return pollData, nil
	}

	questions, err := pollyDb.FindQuestionsByPollId(pollId)
	pollData.Questions = questions
	if err != nil {
		return pollData, err
	}

	options, err := pollyDb.FindOptionsByPollId(pollId)
	pollData.Options = options
	if err != nil {
		return pollData, err
	}

	votes, err := pollyDb.FindVotesByPollId(pollId)
	pollData.Votes = votes
	if err != nil {
		return pollData, err
	}

	creator, err := pollyDb.FindUserById(poll.CreatorId)
	pollData.Creator = creator
	if err != nil {
		return pollData, err
	}

	participants, err := pollyDb.FindParticipantsByPollId(pollId)
	if err != nil {
		return pollData, err
	}

	pollData.Participants = make([]User, len(participants))
	var user User
	for index, participant := range participants {
		user, err = pollyDb.FindUserById(participant.UserId)
		if err != nil {
			return pollData, err
		}

		pollData.Participants[index] = user
	}

	return pollData, nil
}
