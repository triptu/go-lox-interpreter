//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"golox/lox"
	"sync"
	"syscall/js"
)

type runnerState struct {
	isRunning bool
	cancelRun context.CancelFunc
}

var state runnerState

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("loxrun", functionRunner(runLoxCode))

	// this only works right now where is a sleep in the program, where
	// this instruction can actually be receieved
	// coz - goroutines don't actually work in WASM -
	// https://github.com/tinygo-org/tinygo/issues/3095
	// https://github.com/tinygo-org/tinygo/issues/2630
	js.Global().Set("loxstop", functionRunner(stopLoxCode))

	// pause, this is needed for Wasm exports to be visible to JavaScript
	<-c
}

func stopLoxCode(this js.Value, args []js.Value) {
	callbackJs := args[0]
	logToJs(callbackJs, "log", "force stopped manually")
	state.cancelRun()
}

func runLoxCode(this js.Value, args []js.Value) {
	command := args[0].String()
	sourceCode := args[1].String()
	callbackJs := args[2]

	defer func() {
		logToJs(callbackJs, "done", "")
		state.isRunning = false
	}()

	lox.ResetErrorState()

	logOutput := func(s string, isError bool) {
		if isError {
			logToJs(callbackJs, "error", s)
		} else {
			logToJs(callbackJs, "log", s)
		}
	}
	if command != "run" {
		logOutput(fmt.Sprintf("Unknown command: %s", command), true)
		return
	}

	if state.isRunning {
		state.cancelRun()
		logOutput("force stopped, press run again to continue", true)
		return
	}
	state.isRunning = true

	ctx, cancel := context.WithCancel(context.Background())
	state.cancelRun = cancel

	lox.SetLogger(lox.Logger{
		Input: func(prompt string) (string, error) {
			// we'll take input through prompt  in js to keep things simple
			return js.Global().Get("prompt").Invoke(prompt).String(), nil
		},
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
	// panic would also crash coz tinygo doesn't support recover
	// https://github.com/tinygo-org/tinygo/pull/4380
	exitCode := lox.Run([]byte(sourceCode), ctx)
	if exitCode != 0 {
		logOutput(fmt.Sprintf("exit code: %d", exitCode), true)
	}
}

func logToJs(callbackJs js.Value, kind string, msg string) {
	data := make(map[string]interface{})
	data["type"] = kind
	data["data"] = msg
	callbackJs.Invoke(data)
}

func functionRunner(fn func(this js.Value, args []js.Value)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer recover()
			fn(this, args)
		}()
		wg.Wait()
		return nil
	})
}
