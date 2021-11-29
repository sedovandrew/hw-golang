package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	escapeChar = `\`
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packString string) (string, error) {
	var stringBuilder strings.Builder
	var buffer string
	var escaped bool
	for _, char := range packString {
		switch {
		case buffer == escapeChar && !escaped:
			if string(char) == escapeChar {
				escaped = true
			}
			buffer = string(char)
		case buffer == escapeChar && escaped && string(char) == escapeChar:
			stringBuilder.WriteString(buffer)
			escaped = false
		case unicode.IsDigit(char):
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
		default:
			stringBuilder.WriteString(buffer)
			buffer = string(char)
		}
	}
	if buffer == escapeChar && !escaped {
		return "", ErrInvalidString
	}
	stringBuilder.WriteString(buffer)
	return stringBuilder.String(), nil
}
