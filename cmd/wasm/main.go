//go:build js && wasm

package main

import (
	"golox/lox"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("loxrun", js.FuncOf(func(this js.Value, args []js.Value) any {
		callbackJs := args[2]
		callbackGo := func(s string) {
			result := make(map[string]interface{})
			result["type"] = "log"
			result["data"] = s
			callbackJs.Invoke(result)
		}
		runLoxCode(args[0].String(), args[1].String(), callbackGo)

		final := make(map[string]interface{})
		final["type"] = "done"
		final["data"] = ""
		callbackJs.Invoke(final)
		return nil
	}))

	// pause, this is needed for Wasm exports to be visible to JavaScript
	<-c
}

func runLoxCode(command, sourceCode string, printTarget func(string)) {
	if command == "run" {
		lox.Run([]byte(sourceCode), printTarget)
	}
}
