package database

import (
	"polly"

	"gopkg.in/gorp.v1"
)

func (db *Database) AddUser(user *polly.PrivateUser) error {
	return db.dbMap.Insert(user)
}

func (db *Database) AddPoll(poll *polly.Poll) error {
	return db.dbMap.Insert(poll)
}

func (db *Database) AddPollTx(poll *polly.Poll, tx *gorp.Transaction) error {
	return tx.Insert(poll)
}

func (db *Database) AddQuestion(question *polly.Question) error {
	return db.dbMap.Insert(question)
}

func (db *Database) AddQuestionTx(question *polly.Question,
	tx *gorp.Transaction) error {

	return tx.Insert(question)
}

func (db *Database) AddOption(option *polly.Option) error {
	return db.dbMap.Insert(option)
}

func (db *Database) AddOptionTx(option *polly.Option,
	tx *gorp.Transaction) error {

	return tx.Insert(option)
}

func (db *Database) AddVote(vote *polly.Vote) error {
	return db.dbMap.Insert(vote)
}

func AddVoteTx(vote *polly.Vote, tx *gorp.Transaction) error {
	return tx.Insert(vote)
}

func (db *Database) AddParticipantTx(partic *polly.Participant,
	tx *gorp.Transaction) error {

	return tx.Insert(partic)
}

func (db *Database) AddVerToken(verTkn *polly.VerToken) error {
	return db.dbMap.Insert(verTkn)
}
