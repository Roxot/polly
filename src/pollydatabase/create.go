package pollydatabase

func (pollyDb PollyDatabase) AddUser(user *User) error {
	return pollyDb.dbMap.Insert(user)
}

func (pollyDb PollyDatabase) AddPoll(poll *Poll) error {
	return pollyDb.dbMap.Insert(poll)
}

func (pollyDb PollyDatabase) AddQuestion(question *Question) error {
	return pollyDb.dbMap.Insert(question)
}

func (pollyDb PollyDatabase) AddOption(option *Option) error {
	return pollyDb.dbMap.Insert(option)
}

func (pollyDb PollyDatabase) AddVote(vote *Vote) error {
	return pollyDb.dbMap.Insert(vote)
}

func (pollyDb PollyDatabase) AddVerificationToken(vt *VerificationToken) error {
	return pollyDb.dbMap.Insert(vt)
}
