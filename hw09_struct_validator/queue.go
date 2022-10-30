package hw09structvalidator

import (
	"reflect"
	"strings"
)

const (
	validationTagName   = "validate"
	checkSplitSeparator = "|"
	valueSplitSeperator = ":"
)

type checkType struct {
	name  string
	value string
}

type ValidationItem struct {
	rValue       reflect.Value
	rStructField reflect.StructField
	root         bool
}

// checks returns a list of validations
func (vi ValidationItem) checks() ([]checkType, error) {
	t, ok := vi.rStructField.Tag.Lookup(validationTagName)
	if !ok {
		return []checkType{}, nil
	}
	var checks []checkType
	rawChecks := strings.Split(t, checkSplitSeparator)

	for _, rawCheck := range rawChecks {
		check, err := parseCheck(rawCheck)
		if err != nil {
			return []checkType{}, err
		}
		checks = append(checks, check)
	}

	return checks, nil
}

// parseCheck takes a string with validation and returns a checkType object.
func parseCheck(rawCheck string) (checkType, error) {
	nameValueCheck := strings.Split(rawCheck, valueSplitSeperator)
	if nameValueCheck[0] == validateNested {
		return checkType{name: nameValueCheck[0]}, nil
	}
	if len(nameValueCheck) != 2 {
		return checkType{}, ParseCheckError
	}
	return checkType{nameValueCheck[0], nameValueCheck[1]}, nil
}

type ValidationQueue []ValidationItem

// Pop pops an element from the queue.
func (queue *ValidationQueue) Pop() ValidationItem {
	validationItemReflectValue := (*queue)[0]

	copy((*queue)[0:], (*queue)[1:])
	(*queue)[len(*queue)-1] = ValidationItem{}
	*queue = (*queue)[:len(*queue)-1]

	return validationItemReflectValue
}
