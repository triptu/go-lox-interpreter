package lox

import (
	"fmt"
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

type parseError struct {
	line int
	msg  string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("Error at line %d: %s", e.line, e.msg)
}

func newParser[T expr](tokens []token) *parser {
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
func (p *parser) parse() []stmt {
	var statements []stmt
	for !p.isAtEnd() {
		st, err := p.declaration()
		if err == nil {
			statements = append(statements, st)
		} else {
			logError(err.line, err.msg)
			p.consumeCascadingErrors()
		}
	}
	return statements
}

/*
parses single line expression in the code file like - "1+2*3"
*/
func (p *parser) parseExpression() expr {
	expr, err := p.expression()
	if err != nil {
		logError(err.line, err.msg)
	}
	return expr
}

func (p *parser) declaration() (stmt, *parseError) {
	if p.matchIncrement(tVar) {
		return p.varDecl()
	} else if p.matchIncrement(tFun) {
		return p.funDecl("function")
	} else {
		return p.statement()
	}
}

func (p *parser) varDecl() (stmt, *parseError) {
	name, err := p.consumeToken(tIdentifier, "Expected identifier after 'var'")
	if err != nil {
		return nil, err
	}
	var e expr
	if p.matchIncrement(tEqual) {
		e, err = p.expression()
	}
	if err != nil {
		return nil, err
	}
	err = p.eatSemicolon()
	return sVar{
		name:        name,
		initializer: e,
	}, err
}

/*
kind is either "function" or "method"
*/
func (p *parser) funDecl(kind string) (stmt, *parseError) {
	name, err := p.consumeToken(tIdentifier, "expected "+kind+"name")
	if err != nil {
		return nil, err
	}
	err = p.eatToken(tLeftParen, "expected '(' after "+kind+" name")
	if err != nil {
		return nil, err
	}

	var parameters []token
	hasMore := !p.peekMatch(tRightParen)
	for hasMore {
		param, err := p.consumeToken(tIdentifier, "expected parameter name")
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
		hasMore = p.matchIncrement(tComma)
	}
	err = p.eatToken(tRightParen, "expected ')' after parameters")
	if err != nil {
		return nil, err
	}
	err = p.eatToken(tLeftBrace, "expected '{' before "+kind+" body")
	if err != nil {
		return nil, err
	}
	block, err := p.blockRawStmts()
	if err != nil {
		return nil, err
	}
	return sFunction{
		name:       name,
		parameters: parameters,
		body:       block,
	}, nil
}

func (p *parser) statement() (stmt, *parseError) {
	if p.matchIncrement(tPrint) {
		return p.printStmt()
	} else if p.matchIncrement(tLeftBrace) {
		return p.blockStmt()
	} else if p.matchIncrement(tIf) {
		return p.ifStmt()
	} else if p.matchIncrement(tWhile) {
		return p.whileStmt()
	} else if p.matchIncrement(tFor) {
		return p.forStmt()
	} else if p.matchIncrement(tReturn) {
		return p.returnStmt()
	} else {
		return p.exprStmt()
	}
}

func (p *parser) printStmt() (stmt, *parseError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.eatSemicolon()
	return sPrint{
		expression: expr,
	}, err
}

func (p *parser) blockStmt() (stmt, *parseError) {
	statements, err := p.blockRawStmts()
	return sBlock{
		statements: statements,
	}, err
}

func (p *parser) ifStmt() (stmt, *parseError) {
	err := p.eatToken(tLeftParen, "Expected '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.eatToken(tRightParen, "Expected ')' after if condition.")
	if err != nil {
		return nil, err
	}

	ifBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch stmt
	if p.matchIncrement(tElse) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return sIf{
		condition:  condition,
		thenBranch: ifBranch,
		elseBranch: elseBranch,
	}, nil
}

func (p *parser) whileStmt() (stmt, *parseError) {
	err := p.eatToken(tLeftParen, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.eatToken(tRightParen, "Expected ')' after while condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	return sWhile{
		condition: condition,
		body:      body,
	}, err
}

/*
for is implemented in terms of while. a new block is created with initializer
as the first statement. Condition is put in white condition and updater is added
with while's body in a block attached to while's body.
*/
func (p *parser) forStmt() (stmt, *parseError) {
	err := p.eatToken(tLeftParen, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	var initializer stmt
	if !p.matchIncrement(tSemicolon) {
		initializer, err = p.declaration()
		if err != nil {
			return nil, err
		}
	}

	var condition expr
	if p.matchIncrement(tSemicolon) {
		condition = eLiteral{value: true}
	} else {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
		err = p.eatSemicolon()
		if err != nil {
			return nil, err
		}
	}

	var updater expr
	if !p.peekMatch(tRightParen) {
		updater, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.eatToken(tRightParen, "Expected ')' after while condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if updater != nil {
		body = sBlock{[]stmt{body, sExpr{expression: updater}}}
	}
	whileSt := sWhile{
		condition: condition,
		body:      body,
	}

	if initializer == nil {
		return whileSt, nil
	} else {
		return sBlock{[]stmt{initializer, whileSt}}, nil
	}
}

/*
note that just "return ;" is a valid statement, in which case returned
value would be nil
*/
func (p *parser) returnStmt() (stmt, *parseError) {
	var value expr
	var err *parseError
	returnToken := p.tokens[p.curr-1] // the token for the keyword "return"
	if !p.peekMatch(tSemicolon) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	err = p.eatSemicolon()
	return sReturn{
		keyword: returnToken,
		value:   value,
	}, err
}

/*
gives an array of all statements in a block.
Assumes that the "{" has already been consumed.
*/
func (p *parser) blockRawStmts() ([]stmt, *parseError) {
	var statements []stmt
	for !p.isAtEnd() && !p.peekMatch(tRightBrace) {
		st, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, st)
	}
	err := p.eatToken(tRightBrace, "Expected '}' after block")
	return statements, err
}

func (p *parser) exprStmt() (stmt, *parseError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.eatSemicolon()
	return sExpr{
		expression: expr,
	}, err
}

func (p *parser) expression() (expr, *parseError) {
	return p.assignment()
}

/*
as assign is right associative, we use recursion than a loop.
`a = b = c“ should evaluate to `a = (b = c)`
*/
func (p *parser) assignment() (expr, *parseError) {
	expr, err := p.logicOr()
	if err != nil {
		return nil, err
	}

	if p.peekMatch(tEqual) {
		equalsToken := p.tokens[p.curr]
		p.curr++
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}
		varToken, ok := expr.(eVariable)
		if !ok {
			return nil, parseErrorAt(equalsToken, "invalid assignment target.")
		}
		return eAssign{
			name:  varToken.name,
			value: value,
		}, nil
	}

	return expr, nil
}

func (p *parser) logicOr() (expr, *parseError) {
	return p.binaryOp(p.logicAnd, tOr)
}

func (p *parser) logicAnd() (expr, *parseError) {
	return p.binaryOp(p.equality, tAnd)
}

// ==, !=
func (p *parser) equality() (expr, *parseError) {
	return p.binaryOp(p.comparison, tBangEqual, tEqualEqual)
}

// ==, >=, <=, <, >
func (p *parser) comparison() (expr, *parseError) {
	return p.binaryOp(p.term, tGreater, tGreaterEqual, tLess, tLessEqual)
}

func (p *parser) term() (expr, *parseError) {
	return p.binaryOp(p.factor, tPlus, tMinus)
}

func (p *parser) factor() (expr, *parseError) {
	return p.binaryOp(p.unary, tSlash, tStar)
}

/*
nextPrecedenceFn fun is the next precedence level function, which is wrapping the current operator.
for e.g. (4+2)<(3*7), in above the comparison operator is wrapped by term, factor primary on
both sides. the next precedence level for comparison is term.
*/
func (p *parser) binaryOp(nextPrecedenceFn func() (expr, *parseError), tokens ...TokenType) (expr, *parseError) {
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

		if tokens[0] == tOr || tokens[0] == tAnd {
			expr = eLogical{
				left:     expr,
				operator: operator,
				right:    right,
			}
		} else {
			expr = eBinary{
				left:     expr,
				operator: operator,
				right:    right,
			}
		}
	}

	return expr, nil
}

func (p *parser) unary() (expr, *parseError) {
	if p.peekMatch(tBang, tMinus) {
		operator := p.tokens[p.curr]
		p.curr++
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return eUnary{
			operator: operator,
			right:    right,
		}, nil
	}

	return p.call()
}

func (p *parser) call() (expr, *parseError) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.matchIncrement(tLeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *parser) finishCall(callee expr) (expr, *parseError) {
	var arguments []expr
	hasMore := !p.peekMatch(tRightParen)
	for hasMore {
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
		hasMore = p.matchIncrement(tComma)
	}

	paren, err := p.consumeToken(tRightParen, "expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	if len(arguments) >= 255 {
		// just log, not any big error to stop the parsing process itself
		logError(paren.line, "can't have more than 255 arguments to a function.")
	}

	return eCall{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}, nil
}

func (p *parser) primary() (expr, *parseError) {
	token := p.tokens[p.curr]
	p.curr++

	switch token.tokenType {
	case tTrue:
		return eLiteral{value: true}, nil
	case tFalse:
		return eLiteral{value: false}, nil
	case tNil:
		return eLiteral{value: nil}, nil
	case tNumber, tString:
		return eLiteral{value: token.literal}, nil
	case tLeftParen:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		} else if !p.peekMatch(tRightParen) {
			return nil, p.parseErrorPrev("Expected ')' after expression")
		} else {
			p.curr++ // consume the right paren
			return eGrouping{expression: expr}, nil
		}
	case tIdentifier: // variable access
		return eVariable{name: token}, nil
	default:
		errStr := "': Expect expression."
		if arrIncludes(binaryTokens, token.tokenType) {
			errStr = ": Operator found without left-hand operand."
		}
		return nil, parseErrorAt(token, "Error at '"+token.lexeme+errStr)
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

/*
semicolons must be present at the end of every statement, but
we don't actually care about them
*/
func (p *parser) eatSemicolon() *parseError {
	return p.eatToken(tSemicolon, "Expected ';' after expression")
}

/*
read and forget about the token, if it's not what we expect though, return an error
*/
func (p *parser) eatToken(tokenType TokenType, errMsg string) *parseError {
	_, err := p.consumeToken(tokenType, errMsg)
	return err
}

func (p *parser) consumeToken(tokenType TokenType, errMsg string) (token, *parseError) {
	if p.matchIncrement(tokenType) {
		return p.tokens[p.curr-1], nil
	} else {
		return token{}, p.parseErrorPrev(errMsg)
	}
}

func (p *parser) matchIncrement(token TokenType) bool {
	if !p.isAtEnd() && p.tokens[p.curr].tokenType == token {
		p.curr++
		return true
	}
	return false
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

/*
create a parse error at prev token line
*/
func (p *parser) parseErrorPrev(msg string) *parseError {
	return parseErrorAt(p.tokens[p.curr-1], msg)
}

func parseErrorAt(token token, msg string) *parseError {
	return &parseError{
		line: token.line,
		msg:  msg,
	}
}
