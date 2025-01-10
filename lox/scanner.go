package lox

import "fmt"

type scanner struct {
	source string
	tokens []token

	// to keep track of where we're in scanning
	start    int // start index of current lexeme
	curr     int // curr index we're at
	line     int // the line we're at
	lineChar int // the char we're at on the current line

	hasError bool
}

func createScanner(source string) *scanner {
	return &scanner{
		source:   source,
		tokens:   []token{},
		start:    0,
		curr:     0,
		line:     1,
		lineChar: 1,
	}
}

func (s *scanner) scanTokens() []token {
	for s.curr < len(s.source) {
		s.start = s.curr
		s.scanNextToken()
	}

	s.tokens = append(s.tokens, makeEOFToken(s.line, 0))
	return s.tokens
}

func (s *scanner) scanNextToken() {
	c := s.source[s.curr]
	s.curr++
	s.lineChar++
	switch c {
	case '(':
		s.addSimpleToken(tLeftParen)
	case ')':
		s.addSimpleToken(tRightParen)
	case '{':
		s.addSimpleToken(tLeftBrace)
	case '}':
		s.addSimpleToken(tRightBrace)
	case ',':
		s.addSimpleToken(tComma)
	case '.':
		s.addSimpleToken(tDot)
	case '-':
		s.addSimpleToken(tMinus)
	case '+':
		s.addSimpleToken(tPlus)
	case ';':
		s.addSimpleToken(tSemicolon)
	case '/':
		s.addSimpleToken(tSlash)
	case '*':
		s.addSimpleToken(tStar)
	default:
		s.hasError = true
		fmt.Printf("[line %d] Error: Unexpected character: %s\n", s.line, string(c))
	}
}

func (s *scanner) addSimpleToken(tokenType TokenType) {
	s.tokens = append(s.tokens, token{
		tokenType: tokenType,
		lexeme:    s.source[s.start:s.curr],
		literal:   nil,
		line:      s.line,
		column:    s.lineChar,
	})
}

func (s *scanner) addLiteralToken(tokenType TokenType, literal string) {
	s.tokens = append(s.tokens, token{
		tokenType: tokenType,
		lexeme:    s.source[s.start:s.curr],
		literal:   literal,
		line:      s.line,
		column:    s.lineChar,
	})
}
