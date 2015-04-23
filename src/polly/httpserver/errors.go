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
	cBadIdErr          = "Bad id"
	cBadPageErr        = "Bad page"
	cBadEmailErr       = "Bad email"
	cBadDvcTypeErr     = "Bad device type"
	cBadVoteMsgErr     = "Bad vote message"
	cBadVoteTypeErr    = "Bad vote type"
	cNotRegOrBadTknErr = "Not registered or bad token"
	cAccessRightsErr   = "Insufficient access rights"
	cIdListLengthErr   = "Identifier list longer than limit"
	cNoPollErr         = "No such poll"
	cNoUserErr         = "User not found"
	cNoQuestionErr     = "Question not found"
	cNoOptionErr       = "Option not found"
)

/* Handles a bad request. */
func (srv *HTTPServer) handleBadRequest(tag string, msg string, err error,
	writer http.ResponseWriter, req *http.Request) {

	srv.handleErr(tag, msg, fmt.Sprintf(cLogFmt, msg, err), 400, writer, req)
}

/* Handles an illegal operation. */
func (srv *HTTPServer) handleIllegalOperation(tag string, logMsg string,
	writer http.ResponseWriter, req *http.Request) {

	srv.handleErr(tag, "", logMsg, 403, writer, req)
}

/* Handles an authentication error. */
func (srv *HTTPServer) handleAuthError(tag string, err error,
	writer http.ResponseWriter, req *http.Request) {

	writer.Header().Set("WWW-authenticate", "Basic")
	srv.handleErr(tag, "", fmt.Sprintf(cLogFmt, cAuthErr, err), 401,
		writer, req)
}

/* Handles an internal database error. */
func (srv *HTTPServer) handleDatabaseError(tag string, err error,
	writer http.ResponseWriter, req *http.Request) {

	srv.handleErr(tag, "", fmt.Sprintf(cDatabaseErrLog, err), 500, writer, req)
}

/* Handles a marshalling error.  */
func (srv *HTTPServer) handleMarshallingError(tag string, err error,
	writer http.ResponseWriter, req *http.Request) {

	srv.handleErr(tag, "", fmt.Sprintf(cMarshallingErrLog, err), 500,
		writer, req)
}

/* Handles an error occurred while writing the response. */
func (srv *HTTPServer) handleWritingError(tag string, err error,
	writer http.ResponseWriter, req *http.Request) {

	srv.handleErr(tag, "", fmt.Sprintf(cWritingErrLog, err), 500,
		writer, req)
}

/* Generic error handling function. */
func (srv *HTTPServer) handleErr(tag string, msg string, logMsg string,
	statusCode int, writer http.ResponseWriter, req *http.Request) {

	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	srv.logger.Log(tag, logMsg, host)
	http.Error(writer, msg, statusCode)
}
