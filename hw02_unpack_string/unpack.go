package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	escapeChar       = `\`
	doubleEscapeChar = escapeChar + escapeChar
)

var ErrInvalidString = errors.New("invalid string")

// writeCharThroughBuffer writes a character using a buffer.
func writeCharThroughBuffer(stringBuilderPtr *strings.Builder, bufferPtr *string, preChar string, char rune) error {
	if unicode.IsDigit(char) {
		charRepeat, err := strconv.Atoi(string(char))
		if err != nil {
			return err
		}
		stringBuilderPtr.WriteString(strings.Repeat(preChar, charRepeat))
		*bufferPtr = ""
	} else {
		stringBuilderPtr.WriteString(preChar)
		*bufferPtr = string(char)
	}
	return nil
}

func Unpack(packString string) (string, error) {
	var stringBuilder strings.Builder
	var buffer string
	for _, char := range packString {
		switch buffer {
		case "":
			if unicode.IsDigit(char) {
				return "", ErrInvalidString
			}
			buffer = string(char)
		case escapeChar:
			switch {
			case unicode.IsDigit(char):
				buffer = string(char)
			case string(char) == escapeChar:
				buffer = doubleEscapeChar
			default:
				return "", ErrInvalidString
			}
		case doubleEscapeChar:
			err := writeCharThroughBuffer(&stringBuilder, &buffer, escapeChar, char)
			if err != nil {
				return "", err
			}
		default:
			err := writeCharThroughBuffer(&stringBuilder, &buffer, buffer, char)
			if err != nil {
				return "", err
			}
		}
	}

	// Pop buffer
	switch buffer {
	case escapeChar:
		return "", ErrInvalidString
	case doubleEscapeChar:
		stringBuilder.WriteString(escapeChar)
	default:
		stringBuilder.WriteString(buffer)
	}

	return stringBuilder.String(), nil
}
