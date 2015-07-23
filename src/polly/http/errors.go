package httpserver

import (
	"fmt"
	"net"
	"net/http"
)

const (
	cLogFmt            = "%s: %s"
	cDatabaseErrLog    = "Database error: %s"
	cMarshallingErrLog = "Marshalling error: %s"
	cWritingErrLog     = "Error writing response: %s"
	cAuthErr           = "Authentication error"
	cBadPollErr        = "Bad poll"
	cBadJSONErr        = "Bad JSON"
	cBadIDErr          = "Bad id"
	cBadPageErr        = "Bad page"
	cBadEmailErr       = "Bad email"
	cBadDvcTypeErr     = "Bad device type"
	cBadVoteMsgErr     = "Bad vote message"
	cBadVoteTypeErr    = "Bad vote type"
	cNotRegOrBadTknErr = "Not registered or bad token"
	cAccessRightsErr   = "Insufficient access rights"
	cIDListLengthErr   = "IDentifier list longer than limit"
	cNoPollErr         = "No such poll"
	cNoUserErr         = "User not found"
	cNoQuestionErr     = "Question not found"
	cNoOptionErr       = "Option not found"
)

/* Handles a bad requestuest. */
func (server *sServer) handleBadRequest(tag string, msg string, err error,
	writer http.ResponseWriter, request *http.Request) {

	server.handleErr(tag, msg, fmt.Sprintf(cLogFmt, msg, err), 400, writer, request)
}

/* Handles an illegal operation. */
func (server *sServer) handleIllegalOperation(tag string, logMsg string,
	writer http.ResponseWriter, request *http.Request) {

	server.handleErr(tag, "", logMsg, 403, writer, request)
}

/* Handles an authentication error. */
func (server *sServer) handleAuthError(tag string, err error,
	writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("WWW-authenticate", "Basic")
	server.handleErr(tag, "", fmt.Sprintf(cLogFmt, cAuthErr, err), 401,
		writer, request)
}

/* Handles an internal database error. */
func (server *sServer) handleDatabaseError(tag string, err error,
	writer http.ResponseWriter, request *http.Request) {

	server.handleErr(tag, "", fmt.Sprintf(cDatabaseErrLog, err), 500, writer, request)
}

/* Handles a marshalling error.  */
func (server *sServer) handleMarshallingError(tag string, err error,
	writer http.ResponseWriter, request *http.Request) {

	server.handleErr(tag, "", fmt.Sprintf(cMarshallingErrLog, err), 500,
		writer, request)
}

/* Handles an error occurred while writing the response. */
func (server *sServer) handleWritingError(tag string, err error,
	writer http.ResponseWriter, request *http.Request) {

	server.handleErr(tag, "", fmt.Sprintf(cWritingErrLog, err), 500,
		writer, request)
}

/* Generic error handling function. */
func (server *sServer) handleErr(tag string, msg string, logMsg string,
	statusCode int, writer http.ResponseWriter, request *http.Request) {

	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	server.logger.Log(tag, logMsg, host)
	http.Error(writer, msg, statusCode)
}
