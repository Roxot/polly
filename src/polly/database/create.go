package database

func (pollyDb Database) AddUser(user *User) error {
	return pollyDb.dbMap.Insert(user)
}

func (pollyDb Database) AddPoll(poll *Poll) error {
	return pollyDb.dbMap.Insert(poll)
}

func (pollyDb Database) AddQuestion(question *Question) error {
	return pollyDb.dbMap.Insert(question)
}

func (pollyDb Database) AddOption(option *Option) error {
	return pollyDb.dbMap.Insert(option)
}

func (pollyDb Database) AddVote(vote *Vote) error {
	return pollyDb.dbMap.Insert(vote)
}

func (pollyDb Database) AddVerificationToken(vt *VerificationToken) error {
	return pollyDb.dbMap.Insert(vt)
}
