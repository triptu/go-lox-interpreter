package lox

const compileErrorExitCode = 65
const runtimeErrorExitCode = 70

var hasParseError bool
var hasRuntimeError bool

type TokenLogMeta struct {
	Line int
	Col  int
}

// interface as its different for normal run and wasm
type Logger struct {
	Print        func(s string)                       // corresponds to print in lox
	ScanError    func(line int, col int, msg string)  // error during tokenization
	ParseError   func(token TokenLogMeta, msg string) // error during parsing and resolving(static analysis)
	RuntimeError func(token TokenLogMeta, msg string) // error during interpretation
}

var logger Logger

func SetLogger(logger2 Logger) {
	logger = logger2
}

func ResetErrorState() {
	hasParseError = false
	hasRuntimeError = false
}

func logScanError(line int, col int, msg string) {
	hasParseError = true
	logger.ScanError(line, col, msg)
}

func logParseError(token token, msg string) {
	hasParseError = true
	logger.ParseError(TokenLogMeta{Line: token.line, Col: token.column}, msg)
}

/*
this function also panics, as for runtime error we can't proceed further in interpreter
*/
func logRuntimeError(token token, msg string) {
	hasRuntimeError = true
	logger.RuntimeError(TokenLogMeta{Line: token.line, Col: token.column}, msg)
	panic("runtime error")
}
