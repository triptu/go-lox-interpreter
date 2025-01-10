package lox

import "strconv"

type scanner struct {
	source string
	tokens []token

	// to keep track of where we're in scanning
	start    int // start index of current lexeme
	curr     int // curr index we're at
	line     int // the line we're at
	lineChar int // the char we're at on the current line
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
	for !s.isAtEnd() {
		s.start = s.curr
		s.scanNextToken()
	}

	s.tokens = append(s.tokens, makeEOFToken(s.line, 0))
	return s.tokens
}

func (s *scanner) scanNextToken() {
	c := s.source[s.curr]
	s.advance()
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
		if s.peek() == '/' {
			for !s.isAtEnd() && s.peek() != '\n' { // ignore comments
				s.advance()
			}
		} else {
			s.addSimpleToken(tSlash)
		}
	case '*':
		s.addSimpleToken(tStar)
	case ' ', '\t', '\r':
		// ignore whitespace
	case '\n':
		s.line++
		s.lineChar = 1
	case '!':
		s.addConditionalToken(tBang, tBangEqual)
	case '<':
		s.addConditionalToken(tLess, tLessEqual)
	case '>':
		s.addConditionalToken(tGreater, tGreaterEqual)
	case '=':
		s.addConditionalToken(tEqual, tEqualEqual)
	case '"':
		s.scanString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // digits
		s.scanNumber()
	default:
		logError(s.line, string(c))
	}
}

func (s *scanner) advance() {
	s.curr++
	s.lineChar++
}

func (s *scanner) scanString() {
	for !s.isAtEnd() && s.peek() != '"' {
		s.advance()
		if s.peek() == '\n' { // strings can be multiline
			s.line++
		}
	}

	if s.isAtEnd() {
		logError(s.line, "Unterminated string.")
		return
	}

	s.advance() // skip the closing "
	value := s.source[s.start+1 : s.curr-1]
	s.addToken(tString, value)
}

// scan numbers like 1,2, 3.53, etc
func (s *scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()             // consume the "."
		for isDigit(s.peek()) { // consume digits after the decimal
			s.advance()
		}
	}

	// convert string to float
	num_str := s.source[s.start:s.curr]
	num, err := strconv.ParseFloat(num_str, 64)
	if err != nil {
		logError(s.line, err.Error())
		return
	}
	s.addToken(tNumber, num)
}

// these are the characters - !,<,>,= which token type they become depends on if the next character is =
func (s *scanner) addConditionalToken(solo, withEqual TokenType) {
	if s.isAtEnd() || s.source[s.curr] != '=' {
		s.addSimpleToken(solo)
	} else {
		s.advance()
		s.addSimpleToken(withEqual)
	}
}

func (s *scanner) addSimpleToken(tokenType TokenType) {
	s.addToken(tokenType, nil)
}

func (s *scanner) addToken(tokenType TokenType, literal interface{}) {
	s.tokens = append(s.tokens, token{
		tokenType: tokenType,
		lexeme:    s.source[s.start:s.curr],
		literal:   literal,
		line:      s.line,
		column:    s.lineChar,
	})
}

// get the next character safely
func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.curr]
}

func (s *scanner) peekNext() byte {
	if s.curr+1 >= len(s.source) {
		return 0
	}
	return s.source[s.curr+1]
}

func (s *scanner) isAtEnd() bool {
	return s.curr >= len(s.source)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
