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

	inLoop bool

	anonFuncCount int
}

func Parser(tokens []IToken, jffy Jffy) IParser {
	var p IParser = &parser{
		jffy,
		tokens,
		0,
		false,
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
			p.synchronise()
		}
	}()

	if p.match(FUN) {
		return p.function("function")
	}

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *parser) statement() IStmt {
	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(FOR) {
		p.inLoop = true
		s := p.forStatement()
		p.inLoop = false
		return s
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(RETURN) {
		return p.returnStatement()
	}

	if p.match(WHILE) {
		p.inLoop = true
		s := p.whileStatement()
		p.inLoop = false
		return s
	}

	if p.match(BREAK) {
		return p.breakStatement()
	}

	if p.match(CONTINUE) {
		return p.continueStatement()
	}

	if p.match(LEFT_BRACE) {
		return &Block{
			p.block(),
		}
	}

	return p.expressionStatement()
}

func (p *parser) forStatement() IStmt {

	var initialiser IStmt
	if p.match(SEMICOLON) {
		initialiser = nil
	} else if p.match(VAR) {
		initialiser = p.varDeclaration()
	} else {
		initialiser = p.expressionStatement()
	}

	var condition IExpr = nil
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect \";\" after loop condition.")

	var increment IExpr = nil
	if !p.check(LEFT_BRACE) {
		increment = p.expression()
	}

	next := p.peek()
	if next.Type() != LEFT_BRACE {
		p.error(next, "Expect block after condition")
	}

	body := p.statement()

	if increment != nil {
		incExpr := &Expression{
			Expression: increment,
		}

		statements := []IStmt{
			body,
			incExpr,
		}

		body = &Block{
			Statements: statements,
		}
	}

	if condition == nil {
		condition = &Literal{
			true,
		}
	}

	body = &While{
		condition,
		body,
	}

	if initialiser != nil {
		statements := []IStmt{
			initialiser,
			body,
		}

		body = &Block{
			Statements: statements,
		}
	}

	return body
}

func (p *parser) ifStatement() IStmt {
	// Doing statements my way, go like, no () around the condition

	condition := p.expression()

	next := p.peek()
	if next.Type() != LEFT_BRACE {
		p.error(next, "Expect block after if condition.")
	}

	thenBranch := p.statement()
	var elseBranch IStmt = nil

	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &If{
		condition,
		thenBranch,
		elseBranch,
	}
}

func (p *parser) printStatement() IStmt {
	val := p.expression()
	p.consume(SEMICOLON, "Expect \";\" after value.")

	return &Print{
		val,
	}
}

func (p *parser) returnStatement() IStmt {
	keyword := p.previous()

	var val IExpr = nil

	if !p.check(SEMICOLON) {
		val = p.expression()
	}

	p.consume(SEMICOLON, "Expect \";\" after return value.")
	return &Return{
		keyword,
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

func (p *parser) whileStatement() IStmt {
	condition := p.expression()

	next := p.peek()
	if next.Type() != LEFT_BRACE {
		p.error(next, "Expect block after condition")
	}

	body := p.statement()

	return &While{
		condition,
		body,
	}
}

func (p *parser) breakStatement() IStmt {
	if !p.inLoop {
		p.error(p.tokens[p.current], "Expect \"break\" inside loop.")
	}

	p.consume(SEMICOLON, "Expect \";\" after break.")

	return &Break{}
}

func (p *parser) continueStatement() IStmt {
	if !p.inLoop {
		p.error(p.tokens[p.current], "Expect \"continue\" inside loop.")
	}

	p.consume(SEMICOLON, "Expect \";\" after continue.")

	return &Continue{}
}

func (p *parser) expressionStatement() IStmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect \";\" after expression.")

	return &Expression{
		expr,
	}
}

func (p *parser) function(kind string) IStmt {

	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))

	p.consume(LEFT_PAREN, fmt.Sprintf("Expect \"(\" after %s name.", kind))

	params := []IToken{}

	// If there are any parameters
	if !p.check(RIGHT_PAREN) {
		// Always do once, then loop (mimic do while)
		params = append(params, p.consume(IDENTIFIER, "Expect parameter name."))

		for p.match(COMMA) {
			if len(params) >= 255 {
				p.calmError(p.peek(), "Can't have more than 255 parameters")
			}

			params = append(params, p.consume(IDENTIFIER, "Expect parameter name."))
		}
	}

	p.consume(RIGHT_PAREN, "Expect \")\" after parameters.")

	p.consume(LEFT_BRACE, fmt.Sprintf("Expect \" {\" before %s body.", kind))

	body := p.block()

	return &Function{
		name,
		params,
		body,
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

	expr := p.lambda()

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
	}

	return expr

}

func (p *parser) finishLambda() IExpr {
	params := []IToken{}

	if !p.check(RIGHT_PAREN) {
		params = append(params, p.consume(IDENTIFIER, "Expect parameter name."))

		for p.match(COMMA) {
			if len(params) >= 255 {
				p.calmError(p.peek(), "Can't have more than 255 parameters.")
			}

			params = append(params, p.consume(IDENTIFIER, "Expect parameter name."))
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect \")\" after arguments.")

	p.consume(LEFT_BRACE, "Expect \"{\" before anonymous body.")
	body := p.block()

	return &Lambda{
		paren,
		params,
		body,
	}
}

func (p *parser) lambda() IExpr {
	if p.match(FUN) {
		p.consume(LEFT_PAREN, "Expect \"(\" after anonymous function.")
		return p.finishLambda()
	}

	return p.or()
}

func (p *parser) or() IExpr {
	expr := p.and()

	for p.match(OR) {
		op := p.previous()
		right := p.and()
		expr = &Logical{
			expr,
			op,
			right,
		}
	}

	return expr
}

func (p *parser) and() IExpr {
	expr := p.equality()

	for p.match(AND) {
		op := p.previous()
		right := p.equality()
		expr = &Logical{
			expr,
			op,
			right,
		}
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

	return p.call()
}

func (p *parser) finishCall(callee IExpr) IExpr {

	args := []IExpr{}

	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())

		for p.match(COMMA) {
			if len(args) >= 255 {
				p.calmError(p.peek(), "Can't have more than 255 arguments.")
			}
			args = append(args, p.expression())
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect \")\" after arguments.")

	return &Call{
		callee,
		paren,
		args,
	}
}

func (p *parser) call() IExpr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
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
		case BREAK:
		case CONTINUE:
			fmt.Println("DEBUG")
			return

		}

		p.advance()
	}
}
