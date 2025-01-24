package main

import (
	"fmt"
	"os"

	"golox/lox"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		fmt.Fprintln(os.Stderr, "Commands available: tokenize, parse, evaluate, visualize, run")
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
		exitCode := lox.Run(fileContents, func(s string) {
			fmt.Println(s)
		})
		os.Exit(exitCode)
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
