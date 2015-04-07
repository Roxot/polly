package polly

type PrivateUser struct {
	Id          int    `json:"id"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
	Token       string `json:"token"`
	DisplayName string `db:"dislay_name" json:"display_name"`
	DeviceType  int    `db:"device_type" json:"-"`
	DeviceGUID  string `db:"device_guid" json:"-"`
}

type PublicUser struct {
	Id          int    `json:"id`
	PhoneNumber string `json:"phone_number"`
	DisplayName string `json:"display_name"`
}

type Poll struct {
	Id           int
	CreatorId    int   `db:"creator_id"`
	CreationDate int64 `db:"creation_date"`
	Title        string
}

type Question struct {
	Id       int    `json:"id"`
	PollId   int    `db:"poll_id" json:"-"`
	Type     int    `json:"type"`
	Title    string `json:"title"`
	ClientId int    `db:"-" json:"-"`
}

type Option struct {
	Id         int    `json:"id"`
	PollId     int    `db:"poll_id" json:"-"`
	QuestionId int    `db:"question_id" json:"question_id"`
	Type       int    `json:"-"`
	Value      string `json:"value"`
	OptionalId int    `db:"optional_id" json:"-"`
}

type Vote struct {
	Id       int `json:"id"`
	PollId   int `db:"poll_id" json:"-"`
	OptionId int `db:"option_id json:"option_id"`
	UserId   int `db:"user_id" json:"user_id"`
}

type Participant struct {
	Id     int
	UserId int `db:"user_id"`
	PollId int `db:"poll_id"`
}

type VerToken struct {
	Id                int
	PhoneNumber       string `db:"phone_number"`
	VerificationToken string `db:"verification_token"`
}
