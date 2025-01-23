package main

import "golox/lox"

func main() {}

// This function is imported from JavaScript, as it doesn't define a body.
// You should define a function named 'add' in the WebAssembly 'env'
// module from JavaScript.
//
//export add
func add(x, y int) int {
	println("adding %d and %d\n", x, y)
	return x + y
}

// this function is exported to javascript in the wasm module
//
//export wasmLox
func wasmLox(command string, sourceCode string) string {
	if command == "run" {
		lox.Run([]byte(sourceCode))
	}
	return "hello world"
}
