package database

import "polly"

func (db *Database) AddUser(user *polly.PrivateUser) error {
	return db.dbMap.Insert(user)
}

func (db *Database) AddPoll(poll *polly.Poll) error {
	return db.dbMap.Insert(poll)
}

func (db *Database) AddQuestion(question *polly.Question) error {
	return db.dbMap.Insert(question)
}

func (db *Database) AddOption(option *polly.Option) error {
	return db.dbMap.Insert(option)
}

func (db *Database) AddVote(vote *polly.Vote) error {
	return db.dbMap.Insert(vote)
}

func (db *Database) AddVerToken(vt *polly.VerToken) error {
	return db.dbMap.Insert(vt)
}
