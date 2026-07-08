package consts

import "errors"

var (
	ErrorNoRecordFound           = errors.New("No Record Found")
	ErrorUserNotFound            = errors.New("No User Found")
	ErrorSizeExceeded            = errors.New("Request Body Too Large")
	ErrorUNIQUEConstraintFailed  = errors.New("UNIQUE CONSTRAINT failed, PRIMARY KEY repeated")
	ErrorFOREIGNConstraintFailed = errors.New("FOREIGN KEY CONSTRAINT FAILED, inavlid FOREIGN KEY PROVIDED")
)
