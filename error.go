package polly

const (
	baseInternalErr  = 100
	baseIllegalOpErr = 200
	baseBadReqErr    = 300
	baseAuthErr      = 400
)

const (
	ErrCodeDBAdd             = baseInternalErr + iota // 100
	ErrCodeDBGet             = baseInternalErr + iota // 101
	ErrCodeDBUpdate          = baseInternalErr + iota // 102
	ErrCodeDBDelete          = baseInternalErr + iota // 103
	ErrCodeDBBeginTx         = baseInternalErr + iota // 104
	ErrCodeDBCommitTx        = baseInternalErr + iota // 105
	ErrCodeDBSetTxLevel      = baseInternalErr + iota // 106
	ErrCodeMarshal           = baseInternalErr + iota // 107
	ErrCodeDemarshal         = baseInternalErr + iota // 108
	ErrCodeWriteResponse     = baseInternalErr + iota // 109
	ErrCodeCreateHTTPRequest = baseInternalErr + iota // 110
	ErrCodeDoHTTPRequest     = baseInternalErr + iota // 111
	ErrCodeSendNotification  = baseInternalErr + iota // 112
	ErrCodeScheduleClose     = baseInternalErr + iota // 113
	ErrCodeParseInt          = baseInternalErr + iota // 114
)

const (
	ErrCodeIllegalPollAccess = baseIllegalOpErr + iota // 200
	ErrCodeIllegalAddOption  = baseIllegalOpErr + iota // 201
	ErrCodeTooManyIDs        = baseIllegalOpErr + iota // 202
	ErrCodePollClosed        = baseIllegalOpErr + iota // 203
	ErrCodeNotCreator        = baseIllegalOpErr + iota // 204
)

const (
	ErrCodeBadJSON              = baseBadReqErr + iota // 300
	ErrCodeNoUser               = baseBadReqErr + iota // 301
	ErrCodeNoPoll               = baseBadReqErr + iota // 302
	ErrCodeNoQuestion           = baseBadReqErr + iota // 303
	ErrCodeNoOption             = baseBadReqErr + iota // 304
	ErrCodeNoVote               = baseBadReqErr + iota // 305
	ErrCodeBadDeviceType        = baseBadReqErr + iota // 306
	ErrCodeEmptyPoll            = baseBadReqErr + iota // 307
	ErrCodeBadPollType          = baseBadReqErr + iota // 308
	ErrCodeDuplicateParticipant = baseBadReqErr + iota // 309
	ErrCodeBadCreator           = baseBadReqErr + iota // 310
	ErrCodeEmptyQuestion        = baseBadReqErr + iota // 311
	ErrCodeEmptyOption          = baseBadReqErr + iota // 312
	ErrCodeBadVoteType          = baseBadReqErr + iota // 313
	ErrCodeBadPage              = baseBadReqErr + iota // 314
	ErrCodeBadID                = baseBadReqErr + iota // 315
	ErrCodeBadClosingDate       = baseBadReqErr + iota // 316
	ErrCodeNoID                 = baseBadReqErr + iota // 317
	ErrCodeNoDisplayName        = baseBadReqErr + iota // 318
)

const (
	ErrCodeMissingAuth = baseAuthErr + iota // 400
	ErrCodeAuthFailed  = baseAuthErr + iota // 401
	ErrCodeNoFBToken   = baseAuthErr + iota // 402
	ErrCodeBadFBToken  = baseAuthErr + iota // 403
)

var errorMessages = map[int]string{
	ErrCodeDBAdd:             "Failed to add to database.",
	ErrCodeDBGet:             "Failed to retrieve from database.",
	ErrCodeDBUpdate:          "Failed to update database.",
	ErrCodeDBDelete:          "Failed to delete from database.",
	ErrCodeDBBeginTx:         "Failed to start transaction.",
	ErrCodeDBCommitTx:        "Failed to commit transaction.",
	ErrCodeDBSetTxLevel:      "Failed to set transaction level.",
	ErrCodeMarshal:           "Failed to marshal response.",
	ErrCodeDemarshal:         "Failed to demarshal request.",
	ErrCodeWriteResponse:     "Failed to write response.",
	ErrCodeCreateHTTPRequest: "Failed to create HTTP request.",
	ErrCodeDoHTTPRequest:     "Failed to do HTTP request.",
	ErrCodeSendNotification:  "Failed to send notifications.",
	ErrCodeScheduleClose:     "Failed to schedule poll closing event.",
	ErrCodeParseInt:          "Failed to parse integer.",

	ErrCodeIllegalPollAccess: "No access to poll.",
	ErrCodeIllegalAddOption:  "Not allowed to add options.",
	ErrCodeTooManyIDs:        "Too many identifiers provided.",
	ErrCodePollClosed:        "Poll closed.",
	ErrCodeNotCreator:        "No creator access to poll.",

	ErrCodeBadJSON:              "Bad JSON.",
	ErrCodeNoUser:               "No such user.",
	ErrCodeNoPoll:               "No such poll.",
	ErrCodeNoQuestion:           "No such question.",
	ErrCodeNoOption:             "No such option.",
	ErrCodeNoVote:               "No such vote.",
	ErrCodeBadDeviceType:        "Invalid device type.",
	ErrCodeEmptyPoll:            "Empty poll.",
	ErrCodeBadPollType:          "Invalid poll type.",
	ErrCodeDuplicateParticipant: "Duplicate participant.",
	ErrCodeBadCreator:           "Creator not in participants list.",
	ErrCodeEmptyQuestion:        "Empty question.",
	ErrCodeEmptyOption:          "Empty option.",
	ErrCodeBadVoteType:          "Invalid vote type.",
	ErrCodeBadPage:              "Bad page.",
	ErrCodeBadID:                "Bad ID.",
	ErrCodeBadClosingDate:       "Bad closing date.",
	ErrCodeNoID:                 "No ID provided.",
	ErrCodeNoDisplayName:        "No display name provided.",

	ErrCodeMissingAuth: "No authentication provided.",
	ErrCodeAuthFailed:  "Authentication failed.",
	ErrCodeNoFBToken:   "No Facebook token provided.",
	ErrCodeBadFBToken:  "Bad Facebook token.",
}

// Some comment TODO
type Error struct {
	Message string
	Code    int
}

func (err *Error) Error() string {
	return err.Message
}

// Some comment TODO
func NewError(code int) error {
	return &Error{
		Code:    code,
		Message: errorMessages[code],
	}
}
