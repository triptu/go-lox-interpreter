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

func PrintTokens(code []byte) {
	tokens := tokenize(code)
	for _, token := range tokens {
		fmt.Println(token)
	}
	if hasError {
		os.Exit(65)
	}
}

func tokenize(code []byte) []token {
	source := string(code)
	scanner := createScanner(source)
	tokens := scanner.scanTokens()
	return tokens

}

func Parse(code []byte) {
	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("Recovered. Error: ", r)
			if hasError {
				os.Exit(65)
			}
		}
	}()

	tokens := tokenize(code)
	parser := newParser[expr[any]](tokens)
	parsedExpr := parser.expression()
	printer := astPrinter{}
	printer.print(parsedExpr)
}
