//go:build js && wasm

package main

import (
	"fmt"
	"golox/lox"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("loxrun", js.FuncOf(func(this js.Value, args []js.Value) any {
		runLoxCode(args[0].String(), args[1].String(), args[2])
		return nil
	}))

	// pause, this is needed for Wasm exports to be visible to JavaScript
	<-c
}

func runLoxCode(command, sourceCode string, callbackJs js.Value) {
	lox.ResetErrorState()

	logOutput := func(s string, isError bool) {
		result := make(map[string]interface{})
		if isError {
			result["type"] = "error"
		} else {
			result["type"] = "log"
		}
		result["data"] = s
		callbackJs.Invoke(result)
	}

	if command != "run" {
		logOutput(fmt.Sprintf("Unknown command: %s", command), true)
		return
	}

	lox.SetLogger(lox.Logger{
		Print: func(s string) {
			logOutput(s, false)
		},
		ScanError: func(line int, col int, msg string) {
			logOutput(fmt.Sprintf("[line %d:%d] %s", line, col, msg), true)
		},
		ParseError: func(token lox.TokenLogMeta, msg string) {
			logOutput(fmt.Sprintf("[line %d:%d] %s", token.Line, token.Col, msg), true)
		},
		RuntimeError: func(token lox.TokenLogMeta, msg string) {
			logOutput(fmt.Sprintf("[line %d:%d] %s", token.Line, token.Col, msg), true)
		},
	})

	lox.ResetErrorState()
	exitCode := lox.Run([]byte(sourceCode))
	if exitCode != 0 {
		logOutput(fmt.Sprintf("exit code: %d", exitCode), true)
	}

	final := make(map[string]interface{})
	final["type"] = "done"
	final["data"] = ""
	callbackJs.Invoke(final)
}
