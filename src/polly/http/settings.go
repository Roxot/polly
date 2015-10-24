package http

import "time"

const (
	cPollListMax        = 20
	cBulkPollMax        = cPollListMax
	cBulkUserMax        = cBulkPollMax
	cMinPollClosingTime = time.Second * 10
	cMaxPollClosingTime = time.Hour
)
