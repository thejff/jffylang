package jffy

import (
	"fmt"
)

type parseError struct {
	operator IToken
	msg      string
}

type IParser interface {
	Parse() []IStmt
}

type parser struct {
	jffy Jffy

	tokens  []IToken
	current int
}

func Parser(tokens []IToken, jffy Jffy) IParser {
	var p IParser = &parser{
		jffy,
		tokens,
		0,
	}

	return p
}

func (p *parser) Parse() []IStmt {

	// Catch panics
	defer func() {
		if r := recover(); r != nil {
			resp := r.(parseError)
			p.jffy.Error(resp.operator, resp.msg)
		}
	}()

	statements := []IStmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *parser) expression() IExpr {
	return p.assignment()
}

func (p *parser) declaration() IStmt {

	// Catch panics for syncronisation
	defer func() {
		if r := recover(); r != nil {
			// resp := r.(parseError)
			p.synchronise()
		}
	}()

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *parser) statement() IStmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return &Block{
			p.block(),
		}
	}

	return p.expressionStatement()
}

func (p *parser) printStatement() IStmt {
	val := p.expression()
	p.consume(SEMICOLON, "Expect \";\" after value.")

	return &StmtPrint{
		val,
	}
}

func (p *parser) varDeclaration() IStmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initialiser IExpr = nil
	if p.match(EQUAL) {
		initialiser = p.expression()
	}

	p.consume(SEMICOLON, "Expect \";\" after a variable declaration.")

	return &Var{
		name,
		initialiser,
	}
}

func (p *parser) expressionStatement() IStmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect \";\" after expression.")

	return &StmtExpression{
		expr,
	}
}

func (p *parser) block() []IStmt {
	statements := []IStmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect \"}\" after block.")

	return statements
}

func (p *parser) assignment() IExpr {

	expr := p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if v, isVar := expr.(*Variable); isVar {
			name := v.Name
			return &Assign{
				name,
				value,
			}
		}

		p.calmError(equals, "Invalid assignment target.")
		fmt.Println(equals)
	}

	return expr

}

func (p *parser) equality() IExpr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{
			expr,
			operator,
			right,
		}

	}

	return expr
}

func (p *parser) comparison() IExpr {

	var expr IExpr
	expr = p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

		expr = &Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *parser) term() IExpr {
	expr := p.factor()

	for p.match(MINUS, PLUS, DOT_DOT) {
		operator := p.previous()
		right := p.factor()

		expr = &Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *parser) factor() IExpr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()

		expr = &Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *parser) unary() IExpr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()

		return &Unary{
			operator,
			right,
		}
	}

	return p.primary()
}

func (p *parser) primary() IExpr {
	if p.match(FALSE) {
		return &Literal{
			false,
		}
	}

	if p.match(TRUE) {
		return &Literal{
			true,
		}

	}

	if p.match(NIL) {
		return &Literal{
			nil,
		}

	}

	if p.match(NUMBER, STRING) {
		prev := p.previous()
		return &Literal{
			prev.Literal(),
		}

	}

	if p.match(IDENTIFIER) {
		prev := p.previous()
		return &Variable{
			prev,
		}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()

		p.consume(RIGHT_PAREN, "Expect \")\" after expression.")
		return &Grouping{
			expr,
		}
	}

	p.error(p.peek(), "Expect Expression")
	return nil
}

func (p *parser) consume(tType TokenType, msg string) IToken {
	if p.check(tType) {
		return p.advance()
	}

	token := p.peek()
	p.error(token, msg)
	return nil
}

func (p *parser) error(token IToken, msg string) {
	p.jffy.Error(token, msg)

	e := parseError{
		operator: token,
		msg:      msg,
	}

	// Unwind to recover in Parse()
	panic(e)
}

func (p *parser) calmError(token IToken, msg string) {
	p.jffy.Error(token, msg)
}

func (p *parser) match(types ...TokenType) bool {

	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *parser) check(tType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	pVal := p.peek()
	return pVal.Type() == tType
}

func (p *parser) advance() IToken {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *parser) isAtEnd() bool {
	pVal := p.peek()
	return pVal.Type() == EOF
}

func (p *parser) peek() IToken {
	return p.tokens[p.current]
}

func (p *parser) previous() IToken {
	return p.tokens[p.current-1]
}

func (p *parser) synchronise() {
	p.advance()

	for !p.isAtEnd() {
		prvType := p.previous().Type()
		if prvType == SEMICOLON {
			return
		}

		switch p.peek().Type() {

		case CLASS:
		case FOR:
		case FUN:
		case IF:
		case PRINT:
		case RETURN:
		case VAR:
		case WHILE:
			fmt.Println("DEBUG")
			return

		}

		p.advance()
	}
}
