package hw09structvalidator

import "fmt"

var (
	LengthValidationError = fmt.Errorf("The value does not match the specified length")
	InValidationError     = fmt.Errorf("The value is not allowed")
	MinValidationError    = fmt.Errorf("The value is less then the minimum")
	MaxValidationError    = fmt.Errorf("The value is greater then the maximum")
	RegexpValidationError = fmt.Errorf("The value does not match the regular expression")

	NonStructCheckError = fmt.Errorf("The input data is not a structure")
	UnknownCheckError   = fmt.Errorf("Unknown validation error")
	ParseCheckError     = fmt.Errorf("Validation parsing error")
	LengthCheckError    = fmt.Errorf("Length validation error")
	MinCheckError       = fmt.Errorf("Minimum validation error")
	MaxCheckError       = fmt.Errorf("Maximum validation error")
	RegexpCheckError    = fmt.Errorf("Regexp validation error")
)
