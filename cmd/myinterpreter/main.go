package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/lox"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if command == "tokenize" {
		lox.PrintTokens(fileContents)
	} else if command == "parse" {
		lox.Parse(fileContents)
	} else if command == "evaluate" {
		lox.Evaluate(fileContents)
	} else if command == "visualize" {
		lox.Visualize(fileContents)
	} else if command == "run" {
		lox.Run(fileContents)
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
