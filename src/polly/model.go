package polly

const (
	DEVICE_TYPE_ANDROID = 0
	DEVICE_TYPE_IPHONE  = 1

	QUESTION_TYPE_MC   = 0
	QUESTION_TYPE_OPEN = 1

	VOTE_TYPE_NEW    = 0
	VOTE_TYPE_UPVOTE = 1

	EVENT_TYPE_NEW_VOTE = 0
	EVENT_TYPE_UPVOTE   = 1
	EVENT_TYPE_NEW_POLL = 2

	NOTIFICATION_INFO_FIELD = "info"
)

/* Polly primitives */

type PrivateUser struct {
	ID          int64  `json:"id"`
	Token       string `json:"token"`
	DisplayName string `db:"display_name" json:"display_name"`
	DeviceType  int    `db:"device_type" json:"device_type"`
	DeviceGUID  string `db:"device_guid" json:"device_guid"`
}

type Poll struct {
	ID             int64 `json:"poll_id"`
	CreatorID      int64 `db:"creator_id" json:"creator_id"`
	CreationDate   int64 `db:"creation_date" json:"creation_date"`
	LastUpdated    int64 `db:"last_updated" json:"last_updated"`
	SequenceNumber int   `db:"sequence_number" json:"sequence_number"`
}

type Question struct {
	ID     int64  `json:"id"`
	PollID int64  `db:"poll_id" json:"-"`
	Type   int    `json:"type"`
	Title  string `json:"title"`
}

type Option struct {
	ID             int64  `json:"id"`
	PollID         int64  `db:"poll_id" json:"-"`
	QuestionID     int64  `db:"question_id" json:"question_id"`
	Value          string `json:"value"`
	SequenceNumber int    `db:"sequence_number" json:"sequence_number"`
}

type Vote struct {
	ID           int64 `json:"id"`
	PollID       int64 `db:"poll_id" json:"-"`
	OptionID     int64 `db:"option_id" json:"option_id"`
	UserID       int64 `db:"user_id" json:"user_id"`
	CreationDate int64 `db:"creation_date" json:"creation_date"`
}

type Participant struct {
	ID     int64
	UserID int64 `db:"user_id"`
	PollID int64 `db:"poll_id"`
}

/* Partial Polly objects. */

type PublicUser struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
}

type PollSnapshot struct {
	ID             int64 `db:"id" json:"id"`
	LastUpdated    int64 `db:"last_updated" json:"last_updated"`
	SequenceNumber int   `db:"sequence_number" json:"sequence_number"`
}

type DeviceInfo struct {
	DeviceType int    `db:"device_type"`
	DeviceGUID string `db:"device_guid"`
}

/* Polly API message objects */

type PollMessage struct {
	MetaData     Poll         `json:"meta_data"`
	Question     Question     `json:"question"`
	Options      []Option     `json:"options"`
	Votes        []Vote       `json:"votes"`
	Participants []PublicUser `json:"participants"`
}

type PollBulkMessage struct {
	Polls []PollMessage `json:"polls"`
}

type UserBulkMessage struct {
	Users []PublicUser `json:"users"`
}

type VoteMessage struct {
	Type  int    `json:"type"`
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

type VoteResponseMessage struct {
	Option *Option      `json:"option,omitempty"`
	Vote   Vote         `json:"vote"`
	Poll   PollSnapshot `json:"poll"`
}

type UpdateUserMessage struct {
	DeviceGUID  *string `json:"device_guid"`
	DisplayName *string `json:"display_name"`
}

type PollListMessage struct {
	Snapshots  []PollSnapshot `json:"polls"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	NumResults int            `json:"num_results"`
	Total      int64          `json:"total"`
}

type NotificationMessage struct {
	DeviceInfos []DeviceInfo `json:"-"`
	Type        int          `json:"type"`
	User        string       `json:"user"`
	Title       string       `json:"title"`
	PollID      int64        `json:"poll_id"`
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
