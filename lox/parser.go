package lox

import (
	"errors"
)

/*
- The parser makes the AST which can be recursively evaluated by the interpreter.
- Every expression from operators to literal is a node in the AST. The leaf nodes
like literal don't reference anyother node. While non-leaf nodes reference other
nodes.
- The way the tree is built is by calling the methods in order of their precedence. We
  start by calling the lowest precedence method, which in turn calls the next higher
  precdence method, and so on. If the higher predence method actually finds the operator
  it represents, it creates a new node, putting the operator as parent and left/right as the
  expression nodes coming from its recursive down calls.
*/

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
	return p.binaryOp(p.comparison, tBangEqual, tEqualEqual)
}

// ==, >=, <=, <, >
func (p *parser) comparison() (expr[any], error) {
	return p.binaryOp(p.term, tGreater, tGreaterEqual, tLess, tLessEqual)
}

func (p *parser) term() (expr[any], error) {
	return p.binaryOp(p.factor, tPlus, tMinus)
}

func (p *parser) factor() (expr[any], error) {
	return p.binaryOp(p.unary, tSlash, tStar)
}

/*
wrap fun is the next precedence level function, which is wrapping the current operator.
for e.g. (4+2)<(3*7), in above the comparison operator is wrapped by term, factor primary on
both sides. the next precedence level for comparison is term.
*/
func (p *parser) binaryOp(nextPrecedenceFn func() (expr[any], error), tokens ...TokenType) (expr[any], error) {
	expr, err := nextPrecedenceFn()
	if err != nil {
		return nil, err
	}

	// match all operator on the same level
	// notice how this is also making these operators left associative
	// as the newer op encountered on right keeps on becoming a new parent
	for p.peekMatch(tokens...) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := nextPrecedenceFn()
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
		errStr := "': Expect expression."
		if arrIncludes(binaryTokens, token.tokenType) {
			errStr = ": Operator found without left-hand operand."
		}
		logError(token.line, "Error at '"+token.lexeme+errStr)
		p.consumeCascadingErrors()
		return nil, errors.New("unexpected token: " + token.lexeme)
	}
}

func (p *parser) isAtEnd() bool {
	return p.tokens[p.curr].tokenType == tEof
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
	return arrIncludes(tokens, p.tokens[p.curr].tokenType)
}

func arrIncludes[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
