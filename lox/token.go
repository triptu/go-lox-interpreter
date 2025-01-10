package lox

import "fmt"

type TokenType string

const (
	// single character tokens
	tLeftParen  TokenType = "LEFT_PAREN"
	tRightParen           = "RIGHT_PAREN"
	tLeftBrace            = "LEFT_BRACE"
	tRightBrace           = "RIGHT_BRACE"
	tComma                = "COMMA"
	tDot                  = "DOT"
	tMinus                = "MINUS"
	tPlus                 = "PLUS"
	tSemicolon            = "SEMICOLON"
	tSlash                = "SLASH"
	tStar                 = "STAR"

	// conditions(1 or 2 char) tokens
	tBang         = "BANG"
	tBangEqual    = "BANG_EQUAL"
	tEqual        = "EQUAL"
	tEqualEqual   = "EQUAL_EQUAL"
	tGreater      = "GREATER"
	tGreaterEqual = "GREATER_EQUAL"
	tLess         = "LESS"
	tLessEqual    = "LESS_EQUAL"

	// literals
	tIdentifier = "IDENTIFIER"
	tString     = "STRING"
	tNumber     = "NUMBER"

	// keywords
	tAnd    = "AND"
	tClass  = "CLASS"
	tElse   = "ELSE"
	tFalse  = "FALSE"
	tFun    = "FUN"
	tFor    = "FOR"
	tIf     = "IF"
	tNil    = "NIL"
	tOr     = "OR"
	tPrint  = "PRINT"
	tReturn = "RETURN"
	tSuper  = "SUPER"
	tThis   = "THIS"
	tTrue   = "TRUE"
	tVar    = "VAR"
	tWhile  = "WHILE"

	tEof = "EOF"
)

type token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
	column    int
}

func (t token) String() string {
	literalVal := "null"
	if t.tokenType == tNumber {
		f := t.literal.(float64)
		if f == float64(int(f)) {
			literalVal = fmt.Sprintf("%.1f", t.literal) // extra zero for integers
		} else {
			literalVal = fmt.Sprintf("%g", t.literal) // avoid trailing zeroes
		}
	} else if t.literal != nil {
		literalVal = fmt.Sprintf("%v", t.literal)
	}
	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, literalVal)
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
