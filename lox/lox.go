package lox

import (
	"fmt"
	"os"
)

var hasError bool

func logError(line int, msg string) {
	hasError = true
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, msg)
}

func Tokenize(code []byte) {
	source := string(code)
	scanner := createScanner(source)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
	if hasError {
		os.Exit(65)
	}
}

func Parse(code []byte) {
	fmt.Println("parsing code")
	testPrinter()
}
