package lox

import (
	"fmt"
	"os"
)

var hasError bool
var hasRuntimeError bool

func logError(line int, msg string) {
	hasError = true
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, msg)
}

func logRuntimeError(line int, msg string) {
	hasRuntimeError = true
	fmt.Fprintf(os.Stderr, "[line %d] Runtime Error: %s\n", line, msg)
	panic("runtime error")
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
	tokens := tokenize(code)
	parser := newParser[expr[any]](tokens)
	parsedExpr := parser.parse()
	if hasError {
		os.Exit(65)
	} else {
		printer := astPrinter{}
		printer.print(parsedExpr)
	}
}

func Evaluate(code []byte) {
	defer func() {
		if r := recover(); r != nil {
			if !hasRuntimeError {
				fmt.Println("Recovered from run time error panic, Error: ", r)
			}
			os.Exit(70)
		}
	}()

	tokens := tokenize(code)
	if hasError {
		os.Exit(65)
	}

	parser := newParser[expr[any]](tokens)
	parsedExpr := parser.parse()
	if hasError {
		os.Exit(65)
	} else {
		interpreter := interpreter{}
		fmt.Println(getLiteralStr(interpreter.evaluate(parsedExpr)))
		if hasRuntimeError {
			os.Exit(70)
		}
	}
}

func getLiteralStr(literal interface{}) string {
	if literal == nil {
		return "nil"
	}

	switch literal := literal.(type) {
	case float64:
		return fmt.Sprintf("%g", literal) // no extra zeroes
	case string:
		return literal
	default:
		return fmt.Sprintf("%v", literal)
	}
}
