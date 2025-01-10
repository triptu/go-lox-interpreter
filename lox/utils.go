package lox

import "fmt"

func prettyError(line, column int, msg string) {
	fmt.Printf("Error at line %d, column %d: %s\n", line, column, msg)
}
