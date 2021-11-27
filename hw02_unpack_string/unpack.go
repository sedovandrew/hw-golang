package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packString string) (string, error) {
	var stringBuilder strings.Builder
	var buffer string
	for _, char := range packString {
		if unicode.IsDigit(char) {
			charRepeat, err := strconv.Atoi(string(char))
			if err != nil {
				return "", err
			}
			if buffer != "" {
				stringBuilder.WriteString(strings.Repeat(buffer, charRepeat))
				buffer = ""
			} else {
				return "", ErrInvalidString
			}
		} else {
			stringBuilder.WriteString(buffer)
			buffer = string(char)
		}
	}
	stringBuilder.WriteString(buffer)
	return stringBuilder.String(), nil
}
