package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	inSeparator    = ","
	validateNested = "nested"
	validateLen    = "len"
	validateIn     = "in"
	validateMin    = "min"
	validateMax    = "max"
	validateRegexp = "regexp"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	b.WriteString("Found validation errors:\n")
	for _, e := range v {
		b.WriteString(fmt.Sprintf("- in field %q: %q\n", e.Field, e.Err))
	}
	return b.String()
}

// Validate - validate the structure recursively.
func Validate(v interface{}) error {
	validationErrors := ValidationErrors{}

	// Check if not struct
	rValue := reflect.ValueOf(v)
	if rValue.Kind() != reflect.Struct {
		return ErrNonStructCheck
	}

	// Init validation queue
	queue := ValidationQueue{
		ValidationItem{
			rValue: rValue,
			root:   true,
		},
	}

	for len(queue) > 0 {
		err := ValidateItem(&queue, &validationErrors)
		if err != nil {
			return err
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// ValidateItem validates an item from the queue.
func ValidateItem(queue *ValidationQueue, vErrors *ValidationErrors) error {
	validationItemReflectValue := queue.Pop()

	//nolint:exhaustive
	switch validationItemReflectValue.rValue.Kind() {
	case reflect.Struct:
		isValidateStruct, err := IsValidateStruct(validationItemReflectValue)
		if err != nil {
			return err
		}
		// Checking the root structre or with "nested" validation.
		if !validationItemReflectValue.root && !isValidateStruct {
			return nil
		}

		for i := 0; i < validationItemReflectValue.rValue.NumField(); i++ {
			// Skip not public fields
			if !validationItemReflectValue.rValue.Type().Field(i).IsExported() {
				continue
			}

			*queue = append(*queue, ValidationItem{
				rValue:       validationItemReflectValue.rValue.Field(i),
				rStructField: validationItemReflectValue.rValue.Type().Field(i),
			})
		}
	case reflect.String:
		err := ValidateString(validationItemReflectValue, vErrors)
		if err != nil {
			return err
		}
	case reflect.Int:
		err := ValidateInt(validationItemReflectValue, vErrors)
		if err != nil {
			return err
		}
	case reflect.Slice:
		for i := 0; i < validationItemReflectValue.rValue.Len(); i++ {
			*queue = append(*queue, ValidationItem{
				rValue:       validationItemReflectValue.rValue.Index(i),
				rStructField: validationItemReflectValue.rStructField,
			})
		}
	}

	return nil
}

// IsValidateStruct decides if the structure should be validated.
func IsValidateStruct(reflectStruct ValidationItem) (bool, error) {
	checks, err := reflectStruct.checks()
	if err != nil {
		return false, err
	}
	for _, check := range checks {
		if check.name == validateNested {
			return true, nil
		}
	}
	return false, nil
}

// ValidateString validates a field of type string.
func ValidateString(reflectString ValidationItem, vErrors *ValidationErrors) error {
	stringValue := reflectString.rValue.String()

	checks, err := reflectString.checks()
	if err != nil {
		return err
	}
	for _, check := range checks {
		switch check.name {
		case validateLen:
			lengthString, err := strconv.Atoi(check.value)
			if err != nil {
				return ErrLengthCheck
			}
			if lengthString != len(stringValue) {
				*vErrors = append(*vErrors, ValidationError{
					reflectString.rStructField.Name,
					ErrLengthValidation,
				})
			}
		case validateIn:
			allowedValues := strings.Split(check.value, inSeparator)
			if !stringInSlice(stringValue, allowedValues) {
				*vErrors = append(*vErrors, ValidationError{
					reflectString.rStructField.Name,
					ErrInValidation,
				})
			}
		case validateRegexp:
			regexpPattern, err := regexp.Compile(check.value)
			if err != nil {
				return ErrRegexpCheck
			}
			if !regexpPattern.MatchString(stringValue) {
				*vErrors = append(*vErrors, ValidationError{
					reflectString.rStructField.Name,
					ErrRegexpValidation,
				})
			}
		default:
			return ErrUnknownCheck
		}
	}

	return nil
}

// ValidateInt validates a field of type int.
func ValidateInt(reflectInt ValidationItem, vErrors *ValidationErrors) error {
	intValue := reflectInt.rValue.Int()

	checks, err := reflectInt.checks()
	if err != nil {
		return err
	}
	for _, check := range checks {
		switch check.name {
		case validateIn:
			allowedValues := strings.Split(check.value, inSeparator)
			stringValue := strconv.FormatInt(intValue, 10)
			if !stringInSlice(stringValue, allowedValues) {
				*vErrors = append(*vErrors, ValidationError{
					reflectInt.rStructField.Name,
					ErrInValidation,
				})
			}
		case validateMin:
			min, err := strconv.ParseInt(check.value, 10, 64)
			if err != nil {
				return ErrMinCheck
			}
			if intValue < min {
				*vErrors = append(*vErrors, ValidationError{
					reflectInt.rStructField.Name,
					ErrMinValidation,
				})
			}
		case validateMax:
			max, err := strconv.ParseInt(check.value, 10, 64)
			if err != nil {
				return ErrMaxCheck
			}
			if intValue > max {
				*vErrors = append(*vErrors, ValidationError{
					reflectInt.rStructField.Name,
					ErrMaxValidation,
				})
			}
		default:
			return ErrUnknownCheck
		}
	}
	return nil
}
