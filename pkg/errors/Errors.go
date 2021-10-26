package errors

import "fmt"

var (
	ErrContextKeyNotFound       = fmt.Errorf("context key not found")
	ErrContextValueIsNotInteger = fmt.Errorf("context value is not an integer")
	ErrContextValueIsNotFloat   = fmt.Errorf("context value is not a float")
	ErrContextValueIsNotBool    = fmt.Errorf("context value is not a boolean")
	ErrContextValueIsNotTime    = fmt.Errorf("context value is not a time")

	ErrInvalidCronExpression = fmt.Errorf("invalid cron expression")
	ErrEvalError             = fmt.Errorf("error during cel-go evaluation")
	ErrEvalReturnInvalid     = fmt.Errorf("%w : expression evaluation return is invalid", ErrEvalError)
	ErrContextError          = fmt.Errorf("context error")
	ErrStartRequestIfIsFalse = fmt.Errorf("probe request canStart is false")
	ErrCreateHttpClient      = fmt.Errorf("error while creating http client")
	ErrCreateHttpRequest     = fmt.Errorf("error while creating http request")
	ErrHttpCallError         = fmt.Errorf("error while making http call")
	ErrHttpBodyReadError     = fmt.Errorf("error while reading http response body")
	ErrSuccessIfIsFalse      = fmt.Errorf("probe result SuccessIf false")
	ErrFailIfIsTrue          = fmt.Errorf("probe result FailIf true")
)
