package lox

import "fmt"

/**
Tokens are the alphabet of our language grammar. This
file defined all possible tokens and also implements
a pretty print like function for them.
**/

type TokenType int

const (
	// single character tokens
	tLeftParen TokenType = iota
	tRightParen
	tLeftBrace
	tRightBrace
	tComma
	tDot
	tMinus
	tPlus
	tSemicolon
	tSlash
	tStar

	// conditions(1 or 2 char) tokens
	tBang
	tBangEqual
	tEqual
	tEqualEqual
	tGreater
	tGreaterEqual
	tLess
	tLessEqual

	// literals
	tIdentifier
	tString
	tNumber

	// keywords
	tAnd
	tClass
	tElse
	tFalse
	tFun
	tFor
	tIf
	tNil
	tOr
	tPrint
	tReturn
	tSuper
	tThis
	tTrue
	tVar
	tWhile

	tEof
)

var tokenNames = map[TokenType]string{
	tLeftParen:    "LEFT_PAREN",
	tRightParen:   "RIGHT_PAREN",
	tLeftBrace:    "LEFT_BRACE",
	tRightBrace:   "RIGHT_BRACE",
	tComma:        "COMMA",
	tDot:          "DOT",
	tMinus:        "MINUS",
	tPlus:         "PLUS",
	tSemicolon:    "SEMICOLON",
	tSlash:        "SLASH",
	tStar:         "STAR",
	tBang:         "BANG",
	tBangEqual:    "BANG_EQUAL",
	tEqual:        "EQUAL",
	tEqualEqual:   "EQUAL_EQUAL",
	tGreater:      "GREATER",
	tGreaterEqual: "GREATER_EQUAL",
	tLess:         "LESS",
	tLessEqual:    "LESS_EQUAL",
	tIdentifier:   "IDENTIFIER",
	tString:       "STRING",
	tNumber:       "NUMBER",
	tAnd:          "AND",
	tClass:        "CLASS",
	tElse:         "ELSE",
	tFalse:        "FALSE",
	tFun:          "FUN",
	tFor:          "FOR",
	tIf:           "IF",
	tNil:          "NIL",
	tOr:           "OR",
	tPrint:        "PRINT",
	tReturn:       "RETURN",
	tSuper:        "SUPER",
	tThis:         "THIS",
	tTrue:         "TRUE",
	tVar:          "VAR",
	tWhile:        "WHILE",
	tEof:          "EOF",
}

var keywords = map[string]TokenType{
	"and":    tAnd,
	"class":  tClass,
	"else":   tElse,
	"false":  tFalse,
	"for":    tFor,
	"fun":    tFun,
	"if":     tIf,
	"nil":    tNil,
	"or":     tOr,
	"print":  tPrint,
	"return": tReturn,
	"super":  tSuper,
	"this":   tThis,
	"true":   tTrue,
	"var":    tVar,
	"while":  tWhile,
}

type token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
	column    int
}

func (t token) String() string {
	literal := "null"
	if t.literal != nil {
		literal = getTokenLiteralStr(t.literal)
	}
	return fmt.Sprintf("%s %s %s", tokenNames[t.tokenType], t.lexeme, literal)
}

func getTokenLiteralStr(literal interface{}) string {
	if literal == nil {
		return "nil"
	}

	switch literal := literal.(type) {
	case float64:
		if literal == float64(int(literal)) { // is integer
			return fmt.Sprintf("%.1f", literal) // extra zero for integers
		} else {
			return fmt.Sprintf("%g", literal) // avoid trailing zeroes
		}
	case string:
		return literal
	default:
		return fmt.Sprintf("%v", literal)
	}
}

func makeEOFToken(line, column int) token {
	return token{
		tokenType: tEof,
		lexeme:    "",
		literal:   nil,
		line:      line,
		column:    column,
	}
}
