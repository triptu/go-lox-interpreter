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
		callbackJs := args[2]
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
		runLoxCode(args[0].String(), args[1].String(), logOutput)

		final := make(map[string]interface{})
		final["type"] = "done"
		final["data"] = ""
		callbackJs.Invoke(final)
		return nil
	}))

	// pause, this is needed for Wasm exports to be visible to JavaScript
	<-c
}

func runLoxCode(command, sourceCode string, printTarget func(string, bool)) {
	if command == "run" {
		lox.ResetErrorState()
		exitCode := lox.Run([]byte(sourceCode), func(s string) {
			printTarget(s, false)
		})
		if exitCode != 0 {
			printTarget(fmt.Sprintf("exit code: %d", exitCode), true)
		}
	}
}
