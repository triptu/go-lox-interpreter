package lox

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

/*
This is the entry point for Lox exposing public methods for different functionalities.
*/

func PrintTokens(code []byte) {
	tokens := tokenize(code)
	for _, token := range tokens {
		fmt.Println(token)
	}
	if hasParseError {
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
	parser := newParser[expr](tokens)
	parsedExpr := parser.parseExpression()
	if hasParseError {
		os.Exit(65)
	} else {
		printer := astPrinter{}
		printer.print(parsedExpr)
	}
}

func Visualize(code []byte) {
	tokens := tokenize(code)
	parser := newParser[expr](tokens)
	parsedExpr := parser.parseExpression()
	if hasParseError {
		os.Exit(65)
	} else {
		visualizer := NewVisualiseTreeVisitor()
		output_path := "tests/ast_tree"
		if err := visualizer.Visualize(parsedExpr, output_path); err != nil {
			fmt.Printf("Failed to visualize AST: %v\n", err)
		}
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
	if hasParseError {
		os.Exit(65)
	}

	parser := newParser[expr](tokens)
	parsedExpr := parser.parseExpression()
	if hasParseError {
		os.Exit(65)
	} else {
		interpreter := newInterpreter()
		val, _ := interpreter.evaluate(parsedExpr)
		fmt.Println(getLiteralStr(val))
		if hasRuntimeError {
			os.Exit(70)
		}
	}
}

func Run(code []byte, ctx context.Context) (exitCode int) {
	exitCode = 0

	defer func() {
		if r := recover(); r != nil {
			if !hasRuntimeError {
				fmt.Println("Recovered from run time error panic, Error: ", r)
			}
			exitCode = runtimeErrorExitCode
		}
	}()

	tokens := tokenize(code)

	parser := newParser[expr](tokens)
	statements := parser.parse()
	if hasParseError {
		exitCode = compileErrorExitCode
		return
	} else {
		interpreter := newInterpreter()

		resolver := newResolver(interpreter)
		resolver.resolve(statements)
		if hasParseError {
			exitCode = compileErrorExitCode
			return
		}

		interpreter.interpret(statements, ctx)
		if hasRuntimeError {
			exitCode = runtimeErrorExitCode
			return
		}
	}
	return
}

func getLiteralStr(literal interface{}) string {
	if literal == nil {
		return "nil"
	}

	switch literal := literal.(type) {
	case float64:
		return strconv.FormatFloat(literal, 'f', -1, 64)
	case string:
		return literal
	case *loxList:
		return literal.String()
	default:
		return fmt.Sprintf("%v", literal)
	}
}
