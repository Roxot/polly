package polly

import "fmt"

const (
	baseInternalErr  = 100
	baseIllegalOpErr = 200
	baseBadReqErr    = 300
	baseAuthErr      = 400
)

const (
	ErrCodeDatabase          = baseInternalErr + iota // 100
	ErrCodeMarshal                                    // 103
	ErrCodeDemarshal                                  // 104
	ErrCodeWriteResponse                              // 105
	ErrCodeCreateHTTPRequest                          // 106
	ErrCodeDoHTTPRequest                              // 107
	ErrCodeSendNotification                           // 108
	ErrCodeScheduleClose                              // 109
	ErrCodeParseInt                                   // 110
)

const (
	ErrCodeIllegalPollAccess = baseIllegalOpErr + iota // 200
	ErrCodeIllegalAddOption                            // 201
	ErrCodeTooManyIDs                                  // 202
	ErrCodePollClosed                                  // 203
	ErrCodeNotCreator                                  // 204
)

const (
	ErrCodeBadJSON              = baseBadReqErr + iota // 300
	ErrCodeNoUser                                      // 301
	ErrCodeNoPoll                                      // 302
	ErrCodeNoQuestion                                  // 303
	ErrCodeNoOption                                    // 304
	ErrCodeNoVote                                      // 305
	ErrCodeBadDeviceType                               // 306
	ErrCodeEmptyPoll                                   // 307
	ErrCodeBadPollType                                 // 308
	ErrCodeDuplicateParticipant                        // 309
	ErrCodeBadCreator                                  // 310
	ErrCodeEmptyQuestion                               // 311
	ErrCodeEmptyOption                                 // 312
	ErrCodeBadVoteType                                 // 313
	ErrCodeBadPage                                     // 314
	ErrCodeBadID                                       // 315
	ErrCodeBadClosingDate                              // 316
	ErrCodeNoID                                        // 317
	ErrCodeNoDisplayName                               // 318
)

const (
	ErrCodeMissingAuth = baseAuthErr + iota // 400
	ErrCodeAuthFailed                       // 401
	ErrCodeNoFBToken                        // 402
	ErrCodeBadFBToken                       // 403
)

var errorMessages = map[int]string{
	ErrCodeDatabase:          "Internal database error.",
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
	Err     error
	Message string
	Code    int
}

func (err *Error) Error() string {
	return fmt.Sprintf("Polly error %d: %s", err.Code, err.Message)
}

// Some comment TODO
func NewError(err error, code int) error {
	return &Error{
		Err:     err,
		Code:    code,
		Message: errorMessages[code],
	}
}
