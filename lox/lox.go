package lox

import (
	"fmt"
)

func RunCode(code []byte) {
	source := string(code)
	scanner := createScanner(source)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
}
