package polly

const (
	QUESTION_TYPE_MC = 0
	QUESTION_TYPE_OP = 1
	QUESTION_TYPE_DT = 2
)

type PrivateUser struct {
	Id          int    `json:"id"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
	Token       string `json:"token"`
	DisplayName string `db:"dislay_name" json:"display_name"`
	DeviceType  int    `db:"device_type" json:"-"`
	DeviceGUID  string `db:"device_guid" json:"-"`
}

type PublicUser struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	DisplayName string `json:"display_name"`
}

type Poll struct {
	Id           int    `json:"poll_id"`
	CreatorId    int    `db:"creator_id" json:"creator_id"`
	CreationDate int64  `db:"creation_date" json:"creation_date"`
	LastUpdated  int64  `db:"last_updated" json:"last_updated"`
	Title        string `json:"title"`
}

type PollSnapshot struct {
	PollId      int `db:"poll_id" json:"poll_id"`
	LastUpdated int `db:"last_updated" json:"last_updated"`
}

type Question struct {
	Id     int    `json:"id"`
	PollId int    `db:"poll_id" json:"-"`
	Type   int    `json:"type"`
	Title  string `json:"title"`
}

type Option struct {
	Id         int    `json:"id"`
	PollId     int    `db:"poll_id" json:"-"`
	QuestionId int    `db:"question_id" json:"question_id"`
	Value      string `json:"value"`
}

type Vote struct {
	Id       int `json:"id"`
	PollId   int `db:"poll_id" json:"-"`
	OptionId int `db:"option_id" json:"option_id"`
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
