package database

func (db *Database) AddUser(user *PrivateUser) error {
	return db.dbMap.Insert(user)
}

func (db *Database) AddPoll(poll *Poll) error {
	return db.dbMap.Insert(poll)
}

func (db *Database) AddQuestion(question *Question) error {
	return db.dbMap.Insert(question)
}

func (db *Database) AddOption(option *Option) error {
	return db.dbMap.Insert(option)
}

func (db *Database) AddVote(vote *Vote) error {
	return db.dbMap.Insert(vote)
}

func (db *Database) AddVerToken(vt *VerToken) error {
	return db.dbMap.Insert(vt)
}
