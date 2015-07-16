package polly

const (
	DEVICE_TYPE_AD = 0
	DEVICE_TYPE_IP = 1

	QUESTION_TYPE_MC = 0
	QUESTION_TYPE_OP = 1
	QUESTION_TYPE_DT = 2
)

type PrivateUser struct {
	ID          int    `json:"id"`
	Email       string `db:"email" json:"email"`
	Token       string `json:"token"`
	DisplayName string `db:"display_name" json:"display_name"`
	DeviceType  int    `db:"device_type" json:"device_type"`
	DeviceGUID  string `db:"device_guid" json:"device_guid"`
}

type PublicUser struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
}

type Poll struct {
	ID           int   `json:"poll_id"`
	CreatorID    int   `db:"creator_id" json:"creator_id"`
	CreationDate int64 `db:"creation_date" json:"creation_date"`
	LastUpdated  int64 `db:"last_updated" json:"last_updated"`
}

type PollSnapshot struct {
	PollID      int `db:"poll_id" json:"poll_id"`
	LastUpdated int `db:"last_updated" json:"last_updated"`
}

type DeviceInfo struct {
	DeviceType int    `db:"device_type"`
	DeviceGUID string `db:"device_guid"`
}

type Question struct {
	ID     int    `json:"id"`
	PollID int    `db:"poll_id" json:"-"`
	Type   int    `json:"type"`
	Title  string `json:"title"`
}

type Option struct {
	ID         int    `json:"id"`
	PollID     int    `db:"poll_id" json:"-"`
	QuestionID int    `db:"question_id" json:"question_id"`
	Value      string `json:"value"`
}

type Vote struct {
	ID           int   `json:"id"`
	PollID       int   `db:"poll_id" json:"-"`
	OptionID     int   `db:"option_id" json:"option_id"`
	UserID       int   `db:"user_id" json:"user_id"`
	CreationDate int64 `db:"creation_date" json:"creation_date"`
}

type Participant struct {
	ID     int
	UserID int `db:"user_id"`
	PollID int `db:"poll_id"`
}

type VerToken struct {
	ID    int
	Email string `db:"email"`
	Token string `db:"verification_token"`
}
