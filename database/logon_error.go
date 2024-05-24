package database

import "errors"

var BAD_USER_CODE uint8 = 1
var LOCKED_ACCOUNT_CODE uint8 = 2
var FORCED_PASS_RESET_CODE uint8 = 3
var LOGON_COUNT_FAILED_CODE uint8 = 4
var FAILED_TO_SAVE_USER_CODE uint8 = 5

type LogonError struct {
	errorCode uint8
	err       error
}

func (l *LogonError) Error() string {
	return l.err.Error()
}

func (l *LogonError) ErrorCode() int {
	return int(l.errorCode)
}

func LogonErrorNew(err string, code uint8) error {
	var error LogonError
	error.err = errors.New(err)
	error.errorCode = code
	return &error
}
