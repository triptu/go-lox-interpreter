package lox

import (
	"fmt"
	"os"
)

var hasError bool

func logError(line int, msg string) {
	hasError = true
	fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, msg)
}

func RunCode(code []byte) {
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
