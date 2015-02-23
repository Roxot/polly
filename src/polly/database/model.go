package database

const (
	cVerificationTokensTableName = "verification_tokens"
	cUserTableName               = "users"
	cPollTableName               = "polls"
	cQuestionTableName           = "questions"
	cOptionTableName             = "options"
	cVoteTableName               = "votes"
	cParticipantTableName        = "participants"
	cPk                          = "Id"
	cId                          = "id"
	cPhoneNumber                 = "phone_number"
	cToken                       = "token"
	cDisplayName                 = "display_name"
	cDeviceType                  = "device_type"
	cDeviceGUID                  = "device_guid"
	cCreatorId                   = "creator_id"
	cCreationDate                = "creation_date"
	cTitle                       = "title"
	cPollId                      = "poll_id"
	cQuestionId                  = "question_id"
	cValue                       = "value"
	cOptionalId                  = "optional_id"
	cOptionId                    = "option_id"
	cUserId                      = "user_id"
)

type User struct {
	Id          int
	PhoneNumber string `db:"phone_number"`
	Token       string `json:"-"`
	DisplayName string `db:"dislay_name"`
	DeviceType  int    `db:"device_type" json:"-"`
	DeviceGUID  string `db:"device_guid" json:"-"`
}

type Poll struct {
	Id           int
	CreatorId    int   `db:"creator_id"`
	CreationDate int64 `db:"creation_date"`
	Title        string
}

type Question struct {
	Id       int `json:"-"`
	PollId   int `db:"poll_id"`
	Type     int
	Title    string
	ClientId int `db:"-" json:"id"`
}

type Option struct {
	Id         int
	PollId     int `db:"poll_id"`
	QuestionId int `db:"question_id"`
	Type       int
	Value      string
	OptionalId int `db:"optional_id"`
}

type Vote struct {
	Id       int
	PollId   int `db:"poll_id"`
	OptionId int `db:"option_id"`
	UserId   int `db:"user_id"`
}

type Participant struct {
	Id     int
	UserId int `db:"user_id"`
	PollId int `db:"poll_id"`
}

type VerificationToken struct {
	Id                int
	PhoneNumber       string `db:"phone_number"`
	VerificationToken string `db:"verification_token"`
}
