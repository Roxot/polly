package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"polly"
)

type fHeaderHandler func(writer http.ResponseWriter)

const (
	cErrLogFmt         = "%d - \"%s\": %s"
	cUnknownErrLogFmt  = "Unknown error code %d."
	cMarshallingErrFmt = "Error when marshalling error with code %d: %s"
	cWritingErrFmt     = "Error when writing error with code %d: %s."
	cDefaultCode       = 0
	cDefaultMessage    = "Unknown error."
	cDefaultHTTPStatus = http.StatusInternalServerError
	cErrResponseFmt    = "{\n\t\"code\" : %d,\n\t\"message\" : %s\n}"
)

var vDefaultHeaderHandler = setJSONContentTypeHeader

const NO_ERR = 0

const (
	BASE_INT = 100
	BASE_ILL = 200
	BASE_BAD = 300
	BASE_AUT = 400
)

const (
	ERR_INT_DB_ADD       = BASE_INT + iota // 100
	ERR_INT_DB_GET       = BASE_INT + iota // 101
	ERR_INT_DB_UPDATE    = BASE_INT + iota // 102
	ERR_INT_DB_DELETE    = BASE_INT + iota // 103
	ERR_INT_DB_TX_BEGIN  = BASE_INT + iota // 104
	ERR_INT_DB_TX_COMMIT = BASE_INT + iota // 105
	ERR_INT_MARSHALL     = BASE_INT + iota // 106
	ERR_INT_DEMARSHALL   = BASE_INT + iota // 107
	ERR_INT_WRITE        = BASE_INT + iota // 108
	ERR_INT_CREATE_HTTP  = BASE_INT + iota // 109
	ERR_INT_DO_HTTP      = BASE_INT + iota // 110
	ERR_INT_NOTIFICATION = BASE_INT + iota // 111
	ERR_INT_CP_SCHEDULER = BASE_INT + iota // 112
)

const (
	ERR_ILL_POLL_ACCESS  = BASE_ILL + iota // 200
	ERR_ILL_ADD_OPTION   = BASE_ILL + iota // 201
	ERR_ILL_TOO_MANY_IDS = BASE_ILL + iota // 202
	ERR_ILL_POLL_CLOSED  = BASE_ILL + iota // 203
)

const (
	ERR_BAD_JSON                  = BASE_BAD + iota // 300
	ERR_BAD_NO_USER               = BASE_BAD + iota // 301
	ERR_BAD_NO_POLL               = BASE_BAD + iota // 302
	ERR_BAD_NO_QUESTION           = BASE_BAD + iota // 303
	ERR_BAD_NO_OPTION             = BASE_BAD + iota // 304
	ERR_BAD_NO_VOTE               = BASE_BAD + iota // 305
	ERR_BAD_DEVICE_TYPE           = BASE_BAD + iota // 306
	ERR_BAD_EMPTY_POLL            = BASE_BAD + iota // 307
	ERR_BAD_POLL_TYPE             = BASE_BAD + iota // 308
	ERR_BAD_DUPLICATE_PARTICIPANT = BASE_BAD + iota // 309
	ERR_BAD_NO_CREATOR            = BASE_BAD + iota // 310
	ERR_BAD_EMPTY_QUESTION        = BASE_BAD + iota // 311
	ERR_BAD_EMPTY_OPTION          = BASE_BAD + iota // 312
	ERR_BAD_VOTE_TYPE             = BASE_BAD + iota // 313
	ERR_BAD_PAGE                  = BASE_BAD + iota // 314
	ERR_BAD_ID                    = BASE_BAD + iota // 315
	ERR_BAD_CLOSING_DATE          = BASE_BAD + iota // 316
)

const (
	ERR_AUT_NO_AUTH            = BASE_AUT + iota // 400
	ERR_AUT_NO_USER            = BASE_AUT + iota // 401
	ERR_AUT_BAD_TOKEN          = BASE_AUT + iota // 402
	ERR_AUT_NO_OAUTH_HEADERS   = BASE_AUT + iota // 403
	ERR_AUT_BAD_OAUTH_RESPONSE = BASE_AUT + iota // 404
)

var vAPICodeMessages = map[int]string{
	ERR_INT_DB_ADD:       "Failed to add to database.",
	ERR_INT_DB_GET:       "Failed to retrieve from database.",
	ERR_INT_DB_UPDATE:    "Failed to update database.",
	ERR_INT_DB_DELETE:    "Failed to delete from database.",
	ERR_INT_DB_TX_BEGIN:  "Failed to start transaction.",
	ERR_INT_DB_TX_COMMIT: "Failed to commit transaction.",
	ERR_INT_MARSHALL:     "Failed to marshall response.",
	ERR_INT_DEMARSHALL:   "Failed to demarshall request.",
	ERR_INT_WRITE:        "Failed to write response.",
	ERR_INT_CREATE_HTTP:  "Failed to create HTTP request.",
	ERR_INT_DO_HTTP:      "Failed to do HTTP request.",
	ERR_INT_NOTIFICATION: "Failed to send notifications.",
	ERR_INT_CP_SCHEDULER: "Failed to schedule poll closing event.",

	ERR_ILL_POLL_ACCESS:  "No access to poll.",
	ERR_ILL_ADD_OPTION:   "Not allowed to add options.",
	ERR_ILL_TOO_MANY_IDS: "Too many identifiers provided.",
	ERR_ILL_POLL_CLOSED:  "Poll closed.",

	ERR_BAD_JSON:                  "Bad JSON.",
	ERR_BAD_NO_USER:               "No such user.",
	ERR_BAD_NO_POLL:               "No such poll.",
	ERR_BAD_NO_QUESTION:           "No such question.",
	ERR_BAD_NO_OPTION:             "No such option.",
	ERR_BAD_NO_VOTE:               "No such vote.",
	ERR_BAD_DEVICE_TYPE:           "Invalid device type.",
	ERR_BAD_EMPTY_POLL:            "Empty poll.",
	ERR_BAD_POLL_TYPE:             "Invalid poll type.",
	ERR_BAD_DUPLICATE_PARTICIPANT: "Duplicate participant.",
	ERR_BAD_NO_CREATOR:            "Creator not in participants list.",
	ERR_BAD_EMPTY_QUESTION:        "Empty question.",
	ERR_BAD_EMPTY_OPTION:          "Empty option.",
	ERR_BAD_VOTE_TYPE:             "Invalid vote type.",
	ERR_BAD_PAGE:                  "Bad page.",
	ERR_BAD_ID:                    "Bad ID.",
	ERR_BAD_CLOSING_DATE:          "Bad closing date.",

	ERR_AUT_NO_AUTH:            "No authentication provided.",
	ERR_AUT_NO_USER:            "No such user.",
	ERR_AUT_BAD_TOKEN:          "Bad token.",
	ERR_AUT_NO_OAUTH_HEADERS:   "No OAuth Echo headers provided.",
	ERR_AUT_BAD_OAUTH_RESPONSE: "OAuth authentication failed.",
}

var vAPICodeHTTPStatuses = map[int]int{
	ERR_INT_DB_ADD:       http.StatusInternalServerError,
	ERR_INT_DB_GET:       http.StatusInternalServerError,
	ERR_INT_DB_UPDATE:    http.StatusInternalServerError,
	ERR_INT_DB_DELETE:    http.StatusInternalServerError,
	ERR_INT_DB_TX_BEGIN:  http.StatusInternalServerError,
	ERR_INT_DB_TX_COMMIT: http.StatusInternalServerError,
	ERR_INT_MARSHALL:     http.StatusInternalServerError,
	ERR_INT_DEMARSHALL:   http.StatusInternalServerError,
	ERR_INT_WRITE:        http.StatusInternalServerError,
	ERR_INT_CREATE_HTTP:  http.StatusInternalServerError,
	ERR_INT_DO_HTTP:      http.StatusInternalServerError,
	ERR_INT_NOTIFICATION: http.StatusInternalServerError,
	ERR_INT_CP_SCHEDULER: http.StatusInternalServerError,

	ERR_ILL_POLL_ACCESS:  http.StatusForbidden,
	ERR_ILL_ADD_OPTION:   http.StatusForbidden,
	ERR_ILL_TOO_MANY_IDS: http.StatusForbidden,
	ERR_ILL_POLL_CLOSED:  http.StatusForbidden,

	ERR_BAD_JSON:                  http.StatusBadRequest,
	ERR_BAD_NO_USER:               http.StatusBadRequest,
	ERR_BAD_NO_POLL:               http.StatusBadRequest,
	ERR_BAD_NO_QUESTION:           http.StatusBadRequest,
	ERR_BAD_NO_OPTION:             http.StatusBadRequest,
	ERR_BAD_NO_VOTE:               http.StatusBadRequest,
	ERR_BAD_DEVICE_TYPE:           http.StatusBadRequest,
	ERR_BAD_EMPTY_POLL:            http.StatusBadRequest,
	ERR_BAD_POLL_TYPE:             http.StatusBadRequest,
	ERR_BAD_DUPLICATE_PARTICIPANT: http.StatusBadRequest,
	ERR_BAD_NO_CREATOR:            http.StatusBadRequest,
	ERR_BAD_EMPTY_QUESTION:        http.StatusBadRequest,
	ERR_BAD_EMPTY_OPTION:          http.StatusBadRequest,
	ERR_BAD_VOTE_TYPE:             http.StatusBadRequest,
	ERR_BAD_PAGE:                  http.StatusBadRequest,
	ERR_BAD_ID:                    http.StatusBadRequest,
	ERR_BAD_CLOSING_DATE:          http.StatusBadRequest,

	ERR_AUT_NO_AUTH:            http.StatusUnauthorized,
	ERR_AUT_NO_USER:            http.StatusForbidden,
	ERR_AUT_BAD_TOKEN:          http.StatusForbidden,
	ERR_AUT_NO_OAUTH_HEADERS:   http.StatusForbidden,
	ERR_AUT_BAD_OAUTH_RESPONSE: http.StatusForbidden,
}

var vAPICodeHeaderHandler = map[int]fHeaderHandler{ // TODO maybe function fits this better
	ERR_INT_DB_ADD:       setJSONContentTypeHeader,
	ERR_INT_DB_GET:       setJSONContentTypeHeader,
	ERR_INT_DB_UPDATE:    setJSONContentTypeHeader,
	ERR_INT_DB_DELETE:    setJSONContentTypeHeader,
	ERR_INT_DB_TX_BEGIN:  setJSONContentTypeHeader,
	ERR_INT_DB_TX_COMMIT: setJSONContentTypeHeader,
	ERR_INT_MARSHALL:     setJSONContentTypeHeader,
	ERR_INT_DEMARSHALL:   setJSONContentTypeHeader,
	ERR_INT_WRITE:        setJSONContentTypeHeader,
	ERR_INT_CREATE_HTTP:  setJSONContentTypeHeader,
	ERR_INT_DO_HTTP:      setJSONContentTypeHeader,
	ERR_INT_NOTIFICATION: setJSONContentTypeHeader,
	ERR_INT_CP_SCHEDULER: setJSONContentTypeHeader,

	ERR_ILL_POLL_ACCESS:  setJSONContentTypeHeader,
	ERR_ILL_ADD_OPTION:   setJSONContentTypeHeader,
	ERR_ILL_TOO_MANY_IDS: setJSONContentTypeHeader,
	ERR_ILL_POLL_CLOSED:  setJSONContentTypeHeader,

	ERR_BAD_JSON:                  setJSONContentTypeHeader,
	ERR_BAD_NO_USER:               setJSONContentTypeHeader,
	ERR_BAD_NO_POLL:               setJSONContentTypeHeader,
	ERR_BAD_NO_QUESTION:           setJSONContentTypeHeader,
	ERR_BAD_NO_OPTION:             setJSONContentTypeHeader,
	ERR_BAD_NO_VOTE:               setJSONContentTypeHeader,
	ERR_BAD_DEVICE_TYPE:           setJSONContentTypeHeader,
	ERR_BAD_EMPTY_POLL:            setJSONContentTypeHeader,
	ERR_BAD_POLL_TYPE:             setJSONContentTypeHeader,
	ERR_BAD_DUPLICATE_PARTICIPANT: setJSONContentTypeHeader,
	ERR_BAD_NO_CREATOR:            setJSONContentTypeHeader,
	ERR_BAD_EMPTY_QUESTION:        setJSONContentTypeHeader,
	ERR_BAD_EMPTY_OPTION:          setJSONContentTypeHeader,
	ERR_BAD_VOTE_TYPE:             setJSONContentTypeHeader,
	ERR_BAD_PAGE:                  setJSONContentTypeHeader,
	ERR_BAD_ID:                    setJSONContentTypeHeader,
	ERR_BAD_CLOSING_DATE:          setJSONContentTypeHeader,

	ERR_AUT_NO_AUTH:            setAuthenticationChallengeHeaders,
	ERR_AUT_NO_USER:            setJSONContentTypeHeader,
	ERR_AUT_BAD_TOKEN:          setJSONContentTypeHeader,
	ERR_AUT_NO_OAUTH_HEADERS:   setJSONContentTypeHeader,
	ERR_AUT_BAD_OAUTH_RESPONSE: setJSONContentTypeHeader,
}

var vAPICodeShouldLog = map[int]bool{ // TODO maybe function fits this better
	ERR_INT_DB_ADD:       true,
	ERR_INT_DB_GET:       true,
	ERR_INT_DB_UPDATE:    true,
	ERR_INT_DB_DELETE:    true,
	ERR_INT_DB_TX_BEGIN:  true,
	ERR_INT_DB_TX_COMMIT: true,
	ERR_INT_MARSHALL:     true,
	ERR_INT_DEMARSHALL:   true,
	ERR_INT_WRITE:        true,
	ERR_INT_CREATE_HTTP:  true,
	ERR_INT_DO_HTTP:      true,
	ERR_INT_NOTIFICATION: true,
	ERR_INT_CP_SCHEDULER: true,

	ERR_ILL_POLL_ACCESS:  true,
	ERR_ILL_ADD_OPTION:   true,
	ERR_ILL_TOO_MANY_IDS: true,
	ERR_ILL_POLL_CLOSED:  true,

	ERR_BAD_JSON:                  true,
	ERR_BAD_NO_USER:               true,
	ERR_BAD_NO_POLL:               true,
	ERR_BAD_NO_QUESTION:           true,
	ERR_BAD_NO_OPTION:             true,
	ERR_BAD_NO_VOTE:               true,
	ERR_BAD_DEVICE_TYPE:           true,
	ERR_BAD_EMPTY_POLL:            true,
	ERR_BAD_POLL_TYPE:             true,
	ERR_BAD_DUPLICATE_PARTICIPANT: true,
	ERR_BAD_NO_CREATOR:            true,
	ERR_BAD_EMPTY_QUESTION:        true,
	ERR_BAD_EMPTY_OPTION:          true,
	ERR_BAD_VOTE_TYPE:             true,
	ERR_BAD_PAGE:                  true,
	ERR_BAD_ID:                    true,
	ERR_BAD_CLOSING_DATE:          true,

	ERR_AUT_NO_AUTH:            false,
	ERR_AUT_NO_USER:            true,
	ERR_AUT_BAD_TOKEN:          true,
	ERR_AUT_NO_OAUTH_HEADERS:   true,
	ERR_AUT_BAD_OAUTH_RESPONSE: true,
}

func setJSONContentTypeHeader(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func setAuthenticationChallengeHeaders(writer http.ResponseWriter) {
	setJSONContentTypeHeader(writer)
	writer.Header().Set("WWW-authenticate", "Basic")
}

func (server *sServer) respondWithError(errCode int, err error, tag string,
	writer http.ResponseWriter, request *http.Request) {

	knownErr := true
	origin, _, _ := net.SplitHostPort(request.RemoteAddr)

	// check whether we should log
	shouldLog, ok := vAPICodeShouldLog[errCode]
	if !ok {
		knownErr = false
	}

	// get the corresponding error message
	msg, ok := vAPICodeMessages[errCode]
	if !ok {
		msg = cDefaultMessage
		knownErr = false
	}

	// get the correct http status
	httpStatus, ok := vAPICodeHTTPStatuses[errCode]
	if !ok {
		httpStatus = cDefaultHTTPStatus
		knownErr = false
	}

	headerHandler, ok := vAPICodeHeaderHandler[errCode]
	if !ok {
		headerHandler = vDefaultHeaderHandler
		knownErr = false
	}

	// log the error if necessary
	if !knownErr {
		server.logger.Log(tag, fmt.Sprintf(cUnknownErrLogFmt, errCode),
			origin)
	} else if shouldLog {
		server.logger.Log(tag, fmt.Sprintf(cErrLogFmt, errCode, msg, err),
			origin)
	}

	// set the http headers and create the response message
	headerHandler(writer)
	errMsg := polly.ErrorMessage{errCode, msg}

	// marshall the response
	responseBody, err := json.MarshalIndent(errMsg, "", "\t")
	if err != nil {

		// log the error and send the response as plain text
		server.logger.Log(tag, fmt.Sprintf(cMarshallingErrFmt, errCode, err),
			origin)
		http.Error(writer, fmt.Sprintf(cErrResponseFmt, errCode, msg),
			httpStatus)
		return
	}

	// send the response
	writer.WriteHeader(httpStatus)
	_, err = writer.Write(responseBody)
	if err != nil {

		// log the error and send the response as plain text
		server.logger.Log(tag, fmt.Sprintf(cWritingErrFmt, errCode, err),
			origin)
		http.Error(writer, fmt.Sprintf(cErrResponseFmt, errCode, msg),
			httpStatus)
		return
	}
}

func (server *sServer) respondWithJSONBody(writer http.ResponseWriter,
	body []byte) error {

	setJSONContentTypeHeader(writer)
	_, err := writer.Write(body)
	return err
}
