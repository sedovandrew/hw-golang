package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const phrase = "Hello, OTUS!"

func main() {
	reversePhrase := stringutil.Reverse(phrase)
	fmt.Println(reversePhrase)
}
