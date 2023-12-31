package jffy

import (
	"github.com/thejff/jffylang/stack"
)

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
	LAMBDA
)

type LoopType int

const (
	NONLOOP LoopType = iota
	WHILELOOP
	FORLOOP
)

type varState struct {
	token   IToken
	defined bool
	used    bool
}

type resolver struct {
	in       interpreter
	scopes   stack.Stack[[]varState]
	currFunc FunctionType
	currLoop LoopType
}

func Resolver(in interpreter) *resolver {
	scopes := stack.Stack[[]varState]{}
	scopes.Push([]varState{})

	return &resolver{
		in,
		scopes,
		NONE,
		NONLOOP,
	}
}

func (r *resolver) VisitForBlockStmt(stmt *Block) any {

	r.beginScope()
	r.resolve(stmt.Statements)
	r.endScope()

	return nil
}

func (r *resolver) VisitForVarStmt(stmt *Var) any {

	r.declare(stmt.Name)

	if stmt.Initialiser != nil {
		r.resolveExpression(stmt.Initialiser)
	}

	r.define(stmt.Name)

	return nil
}

func (r *resolver) VisitForVariableExpr(expr *Variable) any {
	scope, ok := r.scopes.Pop()
	// if local scope isn't empty, and the variable exists in this scope but is not initialised

	notEmpty := ok
	index := getVarIndex(expr.Name, scope)
	hasEntry := index != -1

	isDefined := false

	if hasEntry {
		isDefined = scope[index].defined
	}

	/* if !isDefined {
		r.in.jffy.Error(expr.Name, "Variable not defined.")
	} */

	if notEmpty && hasEntry && !isDefined {
		r.in.jffy.Error(expr.Name, "Can't read local variable in its own initialiser.")
	}

	if ok && hasEntry {
		scope[index].used = true
		r.scopes.Push(scope)
	}

	r.resolveLocal(expr, expr.Name)

	return nil
}

func (r *resolver) VisitForAssignExpr(expr *Assign) any {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil
}

func (r *resolver) VisitForFunctionStmt(stmt *Function) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)

	return nil
}

func (r *resolver) VisitForExpressionStmt(stmt *Expression) any {
	r.resolveExpression(stmt.Expression)

	return nil
}

func (r *resolver) VisitForIfStmt(stmt *If) any {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.thenBranch)

	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}

	return nil
}

func (r *resolver) VisitForPrintStmt(stmt *Print) any {
	r.resolveExpression(stmt.Expression)

	return nil
}

func (r *resolver) VisitForReturnStmt(stmt *Return) any {

	if r.currFunc == NONE {
		r.in.jffy.Error(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		r.resolveExpression(stmt.Value)
	}

	return nil
}

func (r *resolver) VisitForWhileStmt(stmt *While) any {
	enclosing := r.currLoop
	r.currLoop = WHILELOOP
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.body)

	r.currLoop = enclosing

	return nil
}

func (r *resolver) VisitForBreakStmt(stmt *Break) any {
	if r.currLoop == NONLOOP {
		r.in.jffy.Error(nil, "Can't break outside of loop")
	}
	return nil
}

func (r *resolver) VisitForContinueStmt(stmt *Continue) any {
	if r.currLoop == NONLOOP {
		r.in.jffy.Error(nil, "Can't continue outside of loop")
	}
	return nil
}

func (r *resolver) VisitForBinaryExpr(expr *Binary) any {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)

	return nil
}

func (r *resolver) VisitForLambdaExpr(expr *Lambda) any {

	r.resolveLambda(expr, LAMBDA)

	return nil
}

func (r *resolver) VisitForCallExpr(expr *Call) any {
	r.resolveExpression(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpression(arg)
	}

	return nil
}

func (r *resolver) VisitForGroupingExpr(expr *Grouping) any {
	r.resolveExpression(expr.Expression)

	return nil
}

func (r *resolver) VisitForLiteralExpr(expr *Literal) any {
	return nil
}

func (r *resolver) VisitForLogicalExpr(expr *Logical) any {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)

	return nil
}

func (r *resolver) VisitForUnaryExpr(expr *Unary) any {
	r.resolveExpression(expr.Right)

	return nil
}

func (r *resolver) resolve(stmts []IStmt) {
	for _, s := range stmts {
		r.resolveStatement(s)
	}
}

func (r *resolver) resolveStatement(s IStmt) {
	s.Accept(r)
}

func (r *resolver) resolveExpression(e IExpr) {
	e.Accept(r)
}

func (r *resolver) resolveFunction(function *Function, fnType FunctionType) {

	enclosingFunction := r.currFunc
	r.currFunc = fnType

	r.beginScope()

	for _, p := range function.Params {
		r.declare(p)
		r.define(p)
	}

	r.resolve(function.Body)

	r.endScope()

	r.currFunc = enclosingFunction
}

func (r *resolver) resolveLambda(lambda *Lambda, fnType FunctionType) {
	enclosingFunction := r.currFunc
	r.currFunc = fnType

	r.beginScope()

	for _, p := range lambda.Params {
		r.declare(p)
		r.define(p)
	}

	r.resolve(lambda.Body)

	r.endScope()

	r.currFunc = enclosingFunction
}

func (r *resolver) beginScope() {
	scopeMap := []varState{}
	r.scopes.Push(scopeMap)
}

func (r *resolver) endScope() {
	scope, ok := r.scopes.Pop()
	if !ok {
		// No local variables
		return
	}

	for _, state := range scope {
		if !state.used {
			r.in.jffy.Error(state.token, "Unused variable.")
		}

	}
}

func (r *resolver) declare(name IToken) int {
	scope, ok := r.scopes.Pop()
	// Not ok if stack empty
	if !ok {
		return -1
	}

	for _, v := range scope {
		if v.token.Lexeme() == name.Lexeme() {
			r.in.jffy.Error(name, "There is already a variable with this name in this scope.")
		}
	}

	state := varState{
		token:   name,
		defined: false,
		used:    false,
	}

	scope = append(scope, state)

	r.scopes.Push(scope)

	return len(scope) - 1
}

func (r *resolver) define(name IToken) {
	scope, ok := r.scopes.Pop()
	if !ok {
		return
	}

	index := getVarIndex(name, scope)

	state := scope[index]
	state.defined = true

	scope[index] = state
	r.scopes.Push(scope)
}

func (r *resolver) resolveLocal(expr IExpr, name IToken) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		m, ok := r.scopes.Get(i)
		if ok {
			index := getVarIndex(name, m)
			if index != -1 {
				depth := r.scopes.Size() - 1 - i

				r.in.resolve(expr, depth, index)
				return
			}
		}
	}
}

func getVarIndex(name IToken, scope []varState) int {

	for i, v := range scope {
		if name.Lexeme() == v.token.Lexeme() {
			return i
		}
	}

	return -1
}

func varInScope(name IToken, scope []varState) bool {
	for _, v := range scope {
		if v.token.Lexeme() == name.Lexeme() {
			return true
		}
	}

	return false
}
