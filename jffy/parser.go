package jffy

import (
	"crypto/sha256"
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
		block := p.block()
		uuid := getIdentifier(block)
		return &Block{
			block,
			uuid,
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
		uuid := getIdentifier(true)
		condition = &Literal{
			true,
			uuid,
		}
	}

	uuid := getIdentifier(condition, body)
	body = &While{
		condition,
		body,
		uuid,
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

	uuid := getIdentifier(condition, thenBranch, elseBranch)
	return &If{
		condition,
		thenBranch,
		elseBranch,
		uuid,
	}
}

func (p *parser) printStatement() IStmt {
	val := p.expression()
	p.consume(SEMICOLON, "Expect \";\" after value.")

	uuid := getIdentifier(val)
	return &Print{
		val,
		uuid,
	}
}

func (p *parser) returnStatement() IStmt {
	keyword := p.previous()

	var val IExpr = nil

	if !p.check(SEMICOLON) {
		val = p.expression()
	}

	p.consume(SEMICOLON, "Expect \";\" after return value.")
	uuid := getIdentifier(keyword, val)
	return &Return{
		keyword,
		val,
		uuid,
	}
}

func (p *parser) varDeclaration() IStmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initialiser IExpr = nil
	if p.match(EQUAL) {
		initialiser = p.expression()
	}

	p.consume(SEMICOLON, "Expect \";\" after a variable declaration.")

	uuid := getIdentifier(name, initialiser)
	return &Var{
		name,
		initialiser,
		uuid,
	}
}

func (p *parser) whileStatement() IStmt {
	condition := p.expression()

	next := p.peek()
	if next.Type() != LEFT_BRACE {
		p.error(next, "Expect block after condition")
	}

	body := p.statement()

	uuid := getIdentifier(condition, body)
	return &While{
		condition,
		body,
		uuid,
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

	uuid := getIdentifier(expr)
	return &Expression{
		expr,
		uuid,
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

	uuid := getIdentifier(name, params, body)
	return &Function{
		name,
		params,
		body,
		uuid,
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

			uuid := getIdentifier(name, value)
			return &Assign{
				name,
				value,
				uuid,
			}
		}

		p.calmError(equals, "Invalid assignment target.")
	}

	return expr

}

func getIdentifier(args ...any) string {
	strToHash := ""

	for _, a := range args {
		strToHash += fmt.Sprintf("%p", &a)
	}

	h := sha256.New()
	h.Write([]byte(strToHash))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
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

	uuid := getIdentifier(paren, params, body)

	return &Lambda{
		paren,
		params,
		body,
		uuid,
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

		uuid := getIdentifier(expr, op, right)

		expr = &Logical{
			expr,
			op,
			right,
			uuid,
		}
	}

	return expr
}

func (p *parser) and() IExpr {
	expr := p.equality()

	for p.match(AND) {
		op := p.previous()
		right := p.equality()
		uuid := getIdentifier(expr, op, right)
		expr = &Logical{
			expr,
			op,
			right,
			uuid,
		}
	}

	return expr
}

func (p *parser) equality() IExpr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		uuid := getIdentifier(expr, operator, right)
		expr = &Binary{
			expr,
			operator,
			right,
			uuid,
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

		uuid := getIdentifier(&expr, operator, right)
		expr = &Binary{
			expr,
			operator,
			right,
			uuid,
		}
	}

	return expr
}

func (p *parser) term() IExpr {
	expr := p.factor()

	for p.match(MINUS, PLUS, DOT_DOT) {
		operator := p.previous()
		right := p.factor()

		uuid := getIdentifier(expr, operator, right)
		expr = &Binary{
			expr,
			operator,
			right,
			uuid,
		}
	}

	return expr
}

func (p *parser) factor() IExpr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()

		uuid := getIdentifier(expr, operator, right)
		expr = &Binary{
			expr,
			operator,
			right,
			uuid,
		}
	}

	return expr
}

func (p *parser) unary() IExpr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()

		uuid := getIdentifier(operator, right)
		return &Unary{
			operator,
			right,
			uuid,
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

	uuid := getIdentifier(callee, paren, args)

	return &Call{
		callee,
		paren,
		args,
		uuid,
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
		uuid := getIdentifier(false)
		return &Literal{
			false,
			uuid,
		}
	}

	if p.match(TRUE) {
		uuid := getIdentifier(true)
		return &Literal{
			true,
			uuid,
		}

	}

	if p.match(NIL) {
		uuid := getIdentifier(nil)
		return &Literal{
			nil,
			uuid,
		}

	}

	if p.match(NUMBER, STRING) {
		prev := p.previous()
		uuid := getIdentifier(prev.Literal())
		return &Literal{
			prev.Literal(),
			uuid,
		}

	}

	if p.match(IDENTIFIER) {
		prev := p.previous()
		uuid := getIdentifier(prev)
		return &Variable{
			prev,
			uuid,
		}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()

		p.consume(RIGHT_PAREN, "Expect \")\" after expression.")
		uuid := getIdentifier(expr)
		return &Grouping{
			expr,
			uuid,
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
