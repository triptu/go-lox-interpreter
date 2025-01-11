package lox

type parser[T any] struct {
	tokens []token
	curr   int
}

func newParser[T expr[any]](tokens []token) *parser[T] {
	return &parser[T]{
		tokens: tokens,
		curr:   0,
	}
}

func (p *parser[T]) expression() expr[T] {
	return p.equality()
}

// ==, !=
func (p *parser[T]) equality() expr[T] {
	expr := p.comparison()

	for p.peekMatch(tEqualEqual, tBangEqual) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.comparison()
		expr = eBinary[T]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// ==, >=, <=, <, >
func (p *parser[T]) comparison() expr[T] {
	expr := p.term()

	for p.peekMatch(tGreater, tGreaterEqual, tLess, tLessEqual) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.term()
		expr = eBinary[T]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser[T]) term() expr[T] {
	expr := p.factor()

	for p.peekMatch(tMinus, tPlus) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.factor()
		expr = eBinary[T]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser[T]) factor() expr[T] {
	expr := p.unary()

	for p.peekMatch(tSlash, tStar) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.unary()
		expr = eBinary[T]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser[T]) unary() expr[T] {
	if p.peekMatch(tBang, tMinus) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.unary()
		return eUnary[T]{
			operator: operator,
			right:    right,
		}
	}

	return p.primary()
}

func (p *parser[T]) primary() expr[T] {
	token := p.tokens[p.curr]
	p.curr++

	switch token.tokenType {
	case tTrue:
		return eLiteral[T]{value: true}
	case tFalse:
		return eLiteral[T]{value: false}
	case tNil:
		return eLiteral[T]{value: nil}
	case tNumber, tString:
		return eLiteral[T]{value: token.literal}
	case tLeftParen:
		expr := p.expression()
		if !p.peekMatch(tRightParen) {
			logError(p.tokens[p.curr].line, "Expected ')' after expression")
			panic("Expected ')' after expression")
		}
		p.curr++
		return eGrouping[T]{expression: expr}
	default:
		logError(token.line, "Unexpected token: "+token.lexeme)
		panic("Unexpected token: " + token.lexeme)
	}
}

func (p *parser[T]) isAtEnd() bool {
	return p.curr >= len(p.tokens)
}

// checks if the current token matches any of the given tokens
func (p *parser[T]) peekMatch(tokens ...TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, tokenType := range tokens {
		if p.tokens[p.curr].tokenType == tokenType {
			return true
		}
	}
	return false
}
