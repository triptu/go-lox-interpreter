package lox

import (
	"fmt"
	"os"
)

func RunCode(code []byte) {
	source := string(code)
	scanner := createScanner(source)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
	if scanner.hasError {
		os.Exit(65)
	}
}
