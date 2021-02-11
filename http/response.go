package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/roxot/polly"
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

var vAPICodeHTTPStatuses = map[int]int{
	NO_ERR: http.StatusOK,

	ERR_INT_DB_ADD:             http.StatusInternalServerError,
	ERR_INT_DB_GET:             http.StatusInternalServerError,
	ERR_INT_DB_UPDATE:          http.StatusInternalServerError,
	ERR_INT_DB_DELETE:          http.StatusInternalServerError,
	ERR_INT_DB_TX_BEGIN:        http.StatusInternalServerError,
	ERR_INT_DB_TX_COMMIT:       http.StatusInternalServerError,
	ERR_INT_DB_TX_SET_TX_LEVEL: http.StatusInternalServerError,
	ERR_INT_MARSHALL:           http.StatusInternalServerError,
	ERR_INT_DEMARSHALL:         http.StatusInternalServerError,
	ERR_INT_WRITE:              http.StatusInternalServerError,
	ERR_INT_CREATE_HTTP:        http.StatusInternalServerError,
	ERR_INT_DO_HTTP:            http.StatusInternalServerError,
	ERR_INT_NOTIFICATION:       http.StatusInternalServerError,
	ERR_INT_CP_SCHEDULER:       http.StatusInternalServerError,
	ERR_INT_PARSE_INT:          http.StatusInternalServerError,

	ERR_ILL_POLL_ACCESS:  http.StatusForbidden,
	ERR_ILL_ADD_OPTION:   http.StatusForbidden,
	ERR_ILL_TOO_MANY_IDS: http.StatusForbidden,
	ERR_ILL_POLL_CLOSED:  http.StatusForbidden,
	ERR_ILL_NOT_CREATOR:  http.StatusForbidden,

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
	ERR_BAD_NO_ID:                 http.StatusBadRequest,
	ERR_BAD_NO_DISPLAY_NAME:       http.StatusBadRequest,

	ERR_AUT_NO_AUTH:            http.StatusUnauthorized,
	ERR_AUT_NO_USER:            http.StatusForbidden,
	ERR_AUT_BAD_TOKEN:          http.StatusForbidden,
	ERR_AUT_NO_FACEBOOK_TOKEN:  http.StatusBadRequest,
	ERR_AUT_BAD_FACEBOOK_TOKEN: http.StatusForbidden,
}

var vAPICodeHeaderHandler = map[int]fHeaderHandler{
	NO_ERR: setJSONContentTypeHeader,

	ERR_INT_DB_ADD:             setJSONContentTypeHeader,
	ERR_INT_DB_GET:             setJSONContentTypeHeader,
	ERR_INT_DB_UPDATE:          setJSONContentTypeHeader,
	ERR_INT_DB_DELETE:          setJSONContentTypeHeader,
	ERR_INT_DB_TX_BEGIN:        setJSONContentTypeHeader,
	ERR_INT_DB_TX_COMMIT:       setJSONContentTypeHeader,
	ERR_INT_DB_TX_SET_TX_LEVEL: setJSONContentTypeHeader,
	ERR_INT_MARSHALL:           setJSONContentTypeHeader,
	ERR_INT_DEMARSHALL:         setJSONContentTypeHeader,
	ERR_INT_WRITE:              setJSONContentTypeHeader,
	ERR_INT_CREATE_HTTP:        setJSONContentTypeHeader,
	ERR_INT_DO_HTTP:            setJSONContentTypeHeader,
	ERR_INT_NOTIFICATION:       setJSONContentTypeHeader,
	ERR_INT_CP_SCHEDULER:       setJSONContentTypeHeader,
	ERR_INT_PARSE_INT:          setJSONContentTypeHeader,

	ERR_ILL_POLL_ACCESS:  setJSONContentTypeHeader,
	ERR_ILL_ADD_OPTION:   setJSONContentTypeHeader,
	ERR_ILL_TOO_MANY_IDS: setJSONContentTypeHeader,
	ERR_ILL_POLL_CLOSED:  setJSONContentTypeHeader,
	ERR_ILL_NOT_CREATOR:  setJSONContentTypeHeader,

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
	ERR_BAD_NO_ID:                 setJSONContentTypeHeader,
	ERR_BAD_NO_DISPLAY_NAME:       setJSONContentTypeHeader,

	ERR_AUT_NO_AUTH:            setAuthenticationChallengeHeaders,
	ERR_AUT_NO_USER:            setJSONContentTypeHeader,
	ERR_AUT_BAD_TOKEN:          setJSONContentTypeHeader,
	ERR_AUT_NO_FACEBOOK_TOKEN:  setJSONContentTypeHeader,
	ERR_AUT_BAD_FACEBOOK_TOKEN: setJSONContentTypeHeader,
}

var vAPICodeShouldLog = map[int]bool{
	NO_ERR: false,

	ERR_INT_DB_ADD:             true,
	ERR_INT_DB_GET:             true,
	ERR_INT_DB_UPDATE:          true,
	ERR_INT_DB_DELETE:          true,
	ERR_INT_DB_TX_BEGIN:        true,
	ERR_INT_DB_TX_COMMIT:       true,
	ERR_INT_DB_TX_SET_TX_LEVEL: true,
	ERR_INT_MARSHALL:           true,
	ERR_INT_DEMARSHALL:         true,
	ERR_INT_WRITE:              true,
	ERR_INT_CREATE_HTTP:        true,
	ERR_INT_DO_HTTP:            true,
	ERR_INT_NOTIFICATION:       true,
	ERR_INT_CP_SCHEDULER:       true,
	ERR_INT_PARSE_INT:          true,

	ERR_ILL_POLL_ACCESS:  true,
	ERR_ILL_ADD_OPTION:   true,
	ERR_ILL_TOO_MANY_IDS: true,
	ERR_ILL_POLL_CLOSED:  true,
	ERR_ILL_NOT_CREATOR:  true,

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
	ERR_BAD_NO_ID:                 true,
	ERR_BAD_NO_DISPLAY_NAME:       true,

	ERR_AUT_NO_AUTH:            false,
	ERR_AUT_NO_USER:            true,
	ERR_AUT_BAD_TOKEN:          true,
	ERR_AUT_NO_FACEBOOK_TOKEN:  true,
	ERR_AUT_BAD_FACEBOOK_TOKEN: true,
}

func setJSONContentTypeHeader(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func setAuthenticationChallengeHeaders(writer http.ResponseWriter) {
	setJSONContentTypeHeader(writer)
	writer.Header().Set("WWW-authenticate", "Basic")
}

func (server *sServer) respondOkay(writer http.ResponseWriter,
	request *http.Request) {

	server.respondWithError(NO_ERR, nil, "", writer, request)
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
	errMsg := polly.ErrorMessage{Code: errCode, Message: msg}

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
