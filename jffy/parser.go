package jffy

import "fmt"

type ParseError interface {
	error
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
	expr, err := p.expression()

	if err != nil {
		return nil
	}

	return expr
}

func (p *parser) expression() (IExpr, ParseError) {
	return p.equality()
}

func (p *parser) equality() (IExpr, ParseError) {
	var expr IExpr
	var err ParseError
	expr, err = p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		var right IExpr
		right, err = p.comparison()
		expr = &Binary{
			expr,
			operator,
			right,
		}

	}

	return expr, err
}

func (p *parser) comparison() (IExpr, ParseError) {

	var expr IExpr
	var err ParseError
	expr, err = p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		var right IExpr
		right, err = p.term()
		expr = &Binary{
			expr,
			operator,
			right,
		}
	}

	return expr, err
}

func (p *parser) term() (IExpr, ParseError) {
	var expr IExpr
	var err ParseError

	expr, err = p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()

		var right IExpr
		right, err = p.factor()
		expr = &Binary{
			expr,
			operator,
			right,
		}
	}

	return expr, err
}

func (p *parser) factor() (IExpr, ParseError) {
	var expr IExpr
	var err ParseError
	expr, err = p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()

		var right IExpr
		right, err = p.unary()

		expr = &Binary{
			expr,
			operator,
			right,
		}
	}

	return expr, err
}

func (p *parser) unary() (IExpr, ParseError) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return &Unary{
				operator,
				right,
			},
			err
	}

	return p.primary()
}

func (p *parser) primary() (IExpr, ParseError) {
	if p.match(FALSE) {
		return &Literal{
				false,
			},
			nil
	}

	if p.match(TRUE) {
		return &Literal{
				true,
			},
			nil

	}

	if p.match(NIL) {
		return &Literal{
				nil,
			},
			nil

	}

	if p.match(NUMBER, STRING) {
		prev := p.previous()
		return &Literal{
				prev.Literal(),
			},
			nil

	}

	if p.match(LEFT_PAREN) {
		var expr IExpr
		var err ParseError
		expr, err = p.expression()

		p.consume(RIGHT_PAREN, "Expect \")\" after expression.")
		return &Grouping{
				expr,
			},
			err

	}

	return nil, p.error(p.peek(), "Expect Expression")
}

func (p *parser) consume(tType TokenType, msg string) (IToken, ParseError) {
	if p.check(tType) {
		return p.advance(), nil
	}

	token := p.peek()

	return nil, p.error(token, msg)
}

func (p *parser) error(token IToken, msg string) ParseError {
	p.jffy.Error(token, msg)

	return nil
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
