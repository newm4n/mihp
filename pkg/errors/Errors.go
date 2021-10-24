package errors

import "fmt"

var (
	ErrContextKeyNotFound       = fmt.Errorf("context key not found")
	ErrContextValueIsNotInteger = fmt.Errorf("context value is not an integer")
	ErrContextValueIsNotFloat   = fmt.Errorf("context value is not a float")
	ErrContextValueIsNotBool    = fmt.Errorf("context value is not a boolean")
	ErrContextValueIsNotTime    = fmt.Errorf("context value is not a time")
	ErrInvalidCronExpression    = fmt.Errorf("invalid cron expression")
	ErrEvalReturnInvalid        = fmt.Errorf("expression evaluation return is invalid")
)
