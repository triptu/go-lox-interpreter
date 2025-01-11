package lox

import "errors"

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

/*
parses in a loop, when an error happens, we try to jump to next sane location
in the code, in a way that we avoid cascading errors while still reporting as
many useful errors as possible to the user.
*/
func (p *parser) parse() expr[any] {
	var expr expr[any]
	for !p.isAtEnd() {
		expr, _ = p.expression()
	}
	return expr
}

func (p *parser) expression() (expr[any], error) {
	return p.equality()
}

// ==, !=
func (p *parser) equality() (expr[any], error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.peekMatch(tEqualEqual, tBangEqual) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

// ==, >=, <=, <, >
func (p *parser) comparison() (expr[any], error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.peekMatch(tGreater, tGreaterEqual, tLess, tLessEqual) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (expr[any], error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.peekMatch(tMinus, tPlus) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (expr[any], error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.peekMatch(tSlash, tStar) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = eBinary[any]{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (expr[any], error) {
	if p.peekMatch(tBang, tMinus) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return eUnary[any]{
			operator: operator,
			right:    right,
		}, nil
	}

	return p.primary()
}

func (p *parser) primary() (expr[any], error) {
	token := p.tokens[p.curr]
	p.curr++

	switch token.tokenType {
	case tTrue:
		return eLiteral[any]{value: true}, nil
	case tFalse:
		return eLiteral[any]{value: false}, nil
	case tNil:
		return eLiteral[any]{value: nil}, nil
	case tNumber, tString:
		return eLiteral[any]{value: token.literal}, nil
	case tLeftParen:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		} else if !p.peekMatch(tRightParen) {
			logError(p.tokens[p.curr].line, "Expected ')' after expression")
			p.consumeCascadingErrors()
			return nil, errors.New("expected ')' after expression")
		} else {
			p.curr++ // consume the right paren
			return eGrouping[any]{expression: expr}, nil
		}
	default:
		logError(token.line, "Error at '"+token.lexeme+"': Expect expression.")
		p.consumeCascadingErrors()
		return nil, errors.New("unexpected token: " + token.lexeme)
	}
}

func (p *parser) isAtEnd() bool {
	return p.curr >= len(p.tokens)
}

// when we hit an issue, we increment till we can perhaps restart the parsing process
// this is so we can give the user as much error information as possible
func (p *parser) consumeCascadingErrors() {
	for !p.isAtEnd() {
		if p.peekMatch(tClass, tFun, tVar, tFor, tIf, tWhile, tPrint, tReturn) {
			return
		}
		if p.peekMatch(tSemicolon) {
			p.curr++
			return
		}

		p.curr++
	}
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
