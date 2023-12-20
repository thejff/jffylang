package jffy

import "fmt"

type parseError struct {
	operator IToken
	msg      string
}

type IParser interface {
	Parse() IExpr
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

func (p *parser) Parse() IExpr {

	// Catch panics
	defer func() {
		if r := recover(); r != nil {
			resp := r.(parseError)
			p.jffy.Error(resp.operator, resp.msg)
		}
	}()

	expr := p.expression()

	return expr
}

func (p *parser) expression() IExpr {
	return p.equality()
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

	for p.match(MINUS, PLUS) {
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

func (p *parser) consume(tType TokenType, msg string) {
	if p.check(tType) {
		p.advance()
		return
	}

	token := p.peek()
	p.error(token, msg)
}

func (p *parser) error(token IToken, msg string) {
	p.jffy.Error(token, msg)

	e := &parseError{
		operator: token,
		msg:      msg,
	}

	// Unwind to recover in Parse()
	panic(e)
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
