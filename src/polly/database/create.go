package database

import (
	"polly"

	"gopkg.in/gorp.v1"
)

func (db *Database) AddUser(user *polly.PrivateUser) error {
	return db.mapping.Insert(user)
}

func AddUserTX(user *polly.PrivateUser, tx *gorp.Transaction) error {
	return tx.Insert(user)
}

func (db *Database) AddPoll(poll *polly.Poll) error {
	return db.mapping.Insert(poll)
}

func AddPollTX(poll *polly.Poll, tx *gorp.Transaction) error {
	return tx.Insert(poll)
}

func (db *Database) AddQuestion(question *polly.Question) error {
	return db.mapping.Insert(question)
}

func AddQuestionTX(question *polly.Question, tx *gorp.Transaction) error {
	return tx.Insert(question)
}

func (db *Database) AddOption(option *polly.Option) error {
	return db.mapping.Insert(option)
}

func AddOptionTX(option *polly.Option, tx *gorp.Transaction) error {
	return tx.Insert(option)
}

func (db *Database) AddVote(vote *polly.Vote) error {
	return db.mapping.Insert(vote)
}

func AddVoteTX(vote *polly.Vote, tx *gorp.Transaction) error {
	return tx.Insert(vote)
}

func (db *Database) AddParticipant(participant *polly.Participant) error {
	return db.mapping.Insert(participant)
}

func AddParticipantTX(participant *polly.Participant,
	tx *gorp.Transaction) error {

	return tx.Insert(participant)
}

func (db *Database) AddVerToken(verToken *polly.VerToken) error {
	return db.mapping.Insert(verToken)
}

func AddVerTokenTX(verToken *polly.VerToken, tx *gorp.Transaction) error {
	return tx.Insert(verToken)
}
