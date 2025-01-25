package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	lox.SetLogger(lox.Logger{
		Input: func(prompt string) (string, error) {
			fmt.Print(prompt)
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return "", err
			}
			return input, nil
		},
		Print: func(s string) {
			fmt.Println(s)
		},
		ScanError: func(line int, col int, msg string) {
			fmt.Fprintf(os.Stderr, "[line %d:%d] %s\n", line, col, msg)
		},
		ParseError: func(token lox.TokenLogMeta, msg string) {
			fmt.Fprintf(os.Stderr, "[line %d:%d] %s\n", token.Line, token.Col, msg)
		},
		RuntimeError: func(token lox.TokenLogMeta, msg string) {
			fmt.Fprintf(os.Stderr, "%s\n", msg)
			fmt.Fprintf(os.Stderr, "[line %d:%d] %s\n", token.Line, token.Col, msg)
		},
	})

	if command == "tokenize" {
		lox.PrintTokens(fileContents)
	} else if command == "parse" {
		lox.Parse(fileContents)
	} else if command == "evaluate" {
		lox.Evaluate(fileContents)
	} else if command == "visualize" {
		lox.Visualize(fileContents)
	} else if command == "run" {
		// Create context that listens for the interrupt signal from the OS for graceful stop
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		exitCode := lox.Run(fileContents, ctx)
		os.Exit(exitCode)
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
