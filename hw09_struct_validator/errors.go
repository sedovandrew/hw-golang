package hw09structvalidator

import "errors"

var (
	ErrLengthValidation = errors.New("the value does not match the specified length")
	ErrInValidation     = errors.New("the value is not allowed")
	ErrMinValidation    = errors.New("the value is less then the minimum")
	ErrMaxValidation    = errors.New("the value is greater then the maximum")
	ErrRegexpValidation = errors.New("the value does not match the regular expression")

	ErrNonStructCheck = errors.New("the input data is not a structure")
	ErrUnknownCheck   = errors.New("unknown validation error")
	ErrParseCheck     = errors.New("validation parsing error")
	ErrLengthCheck    = errors.New("length validation error")
	ErrMinCheck       = errors.New("minimum validation error")
	ErrMaxCheck       = errors.New("maximum validation error")
	ErrRegexpCheck    = errors.New("regexp validation error")
)
