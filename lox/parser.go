package lox

type parser struct {
	tokens []token
	curr   int
}

func newParser[T expr[any]](tokens []token) *parser {
	return &parser{
		tokens: tokens,
		curr:   0,
	}
}

func (p *parser) expression() expr[any] {
	return p.equality()
}

// ==, !=
func (p *parser) equality() expr[any] {
	expr := p.comparison()

	for p.peekMatch(tEqualEqual, tBangEqual) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.comparison()
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// ==, >=, <=, <, >
func (p *parser) comparison() expr[any] {
	expr := p.term()

	for p.peekMatch(tGreater, tGreaterEqual, tLess, tLessEqual) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.term()
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) term() expr[any] {
	expr := p.factor()

	for p.peekMatch(tMinus, tPlus) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.factor()
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) factor() expr[any] {
	expr := p.unary()

	for p.peekMatch(tSlash, tStar) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.unary()
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) unary() expr[any] {
	if p.peekMatch(tBang, tMinus) {
		operator := p.tokens[p.curr]
		p.curr++
		right := p.unary()
		return eUnary[any]{
			operator: operator,
			right:    right,
		}
	}

	return p.primary()
}

func (p *parser) primary() expr[any] {
	token := p.tokens[p.curr]
	p.curr++

	switch token.tokenType {
	case tTrue:
		return eLiteral[any]{value: true}
	case tFalse:
		return eLiteral[any]{value: false}
	case tNil:
		return eLiteral[any]{value: nil}
	case tNumber, tString:
		return eLiteral[any]{value: token.literal}
	case tLeftParen:
		expr := p.expression()
		if !p.peekMatch(tRightParen) {
			logError(p.tokens[p.curr].line, "Expected ')' after expression")
			panic("Expected ')' after expression")
		}
		p.curr++
		return eGrouping[any]{expression: expr}
	default:
		logError(token.line, "Unexpected token: "+token.lexeme)
		panic("Unexpected token: " + token.lexeme)
	}
}

func (p *parser) isAtEnd() bool {
	return p.curr >= len(p.tokens)
}

// checks if the current token matches any of the given tokens
func (p *parser) peekMatch(tokens ...TokenType) bool {
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
