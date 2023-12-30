package jffy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LoopControl struct {
	doBreak    bool
	doContinue bool
}

type runtimeError struct {
	operator IToken
	msg      string
}

type dist struct {
	distance int
	set      bool
}

type interpreter struct {
	jffy     Jffy
	env      Environment
	globals  Environment
	locals   map[string]dist
	loopCtrl LoopControl
}

func Interpreter(jffy Jffy) interpreter {
	env := GlobalEnv()

	env.Define("clock", &Clock{})

	l := LoopControl{
		false,
		false,
	}

	locals := make(map[string]dist)

	return interpreter{
		jffy:     jffy,
		env:      env,
		globals:  env,
		locals:   locals,
		loopCtrl: l,
	}

}

func (in *interpreter) Interpret(stmts []IStmt, jffy Jffy) {
	in.jffy = jffy

	// Catch panics
	defer func() {
		if r := recover(); r != nil {
			resp, ok := r.(runtimeError)
			if ok {
				/* fmt.Printf("OP type: %s\n",
					resp.msg,
				) */
				in.jffy.RuntimeError(resp.operator, resp.msg)
			} else {
				err, ok := r.(error)
				if ok {
					in.jffy.RuntimeError(nil, err.Error())
				} else {
					in.jffy.RuntimeError(nil, "Unknown error")
				}

			}
		}
	}()

	for _, s := range stmts {
		in.execute(s)
	}

}

func (in *interpreter) VisitForBinaryExpr(b *Binary) any {
	left := in.evaluate(b.Left)
	right := in.evaluate(b.Right)

	// Non mathematical operations
	switch b.Operator.Type() {

	case DOT_DOT:
		// concatenateOperands will panic if these aren't strings
		return in.concatenateOperands(b.Operator, left, right)

	case BANG_EQUAL:
		return !isEqual(left, right)

	case EQUAL_EQUAL:
		return isEqual(left, right)
	}

	// panics if one not a number
	in.checkNumberOperands(b.Operator, left, right)

	// All below cases require the numbers to be floats
	lVal, rVal, err := lrToFloat(left, right)
	if err != nil {
		in.handleError(b.Operator, err.Error())
		return nil
	}

	// Mathematical operations
	switch b.Operator.Type() {

	case GREATER:
		return lVal > rVal

	case GREATER_EQUAL:
		return lVal >= rVal

	case LESS:
		return lVal < rVal

	case LESS_EQUAL:
		return lVal <= rVal

	case MINUS:
		return lVal - rVal

	case SLASH:
		return lVal / rVal

	case STAR:
		return lVal * rVal

	case PLUS:
		return lVal + rVal

	}

	return nil
}

func (in *interpreter) VisitForLambdaExpr(l *Lambda) any {
	return NewAnonymousFunction(l, in.env)
}

func (in *interpreter) VisitForCallExpr(c *Call) any {
	callee := in.evaluate(c.Callee)

	args := []any{}

	for _, arg := range c.Arguments {
		args = append(args, in.evaluate(arg))
	}

	fun, isCallee := callee.(ICallable)
	if !isCallee {
		in.handleError(c.Paren, "Can only call functions and classes.")
	}

	if len(args) != fun.Arity() {
		in.handleError(
			c.Paren,
			fmt.Sprintf(
				"Expected %d arguments but got %d.",
				fun.Arity(),
				len(args),
			),
		)
	}

	return fun.Call(in, args)
}

func (in *interpreter) VisitForGroupingExpr(g *Grouping) any {
	return in.evaluate(g.Expression)
}

func (in *interpreter) VisitForLiteralExpr(l *Literal) any {
	return l.Value
}

func (in *interpreter) VisitForLogicalExpr(l *Logical) any {
	left := in.evaluate(l.Left)

	if l.Operator.Type() == OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return in.evaluate(l.Right)
}

func (in *interpreter) VisitForUnaryExpr(u *Unary) any {
	right := in.evaluate(u.Right)

	switch u.Operator.Type() {

	case MINUS:
		if isValid := in.checkNumberOperand(u.Operator, right); !isValid {
			return nil
		}

		f, err := anyToFloat(right)
		if err != nil {
			in.handleError(u.Operator, err.Error())
		}

		return -f

	case BANG:
		return !isTruthy(right)

	}

	return nil
}

func (in *interpreter) evaluate(expr IExpr) any {
	return expr.Accept(in)
}

func (in *interpreter) execute(stmt IStmt) any {
	if !in.loopCtrl.doContinue {
		stmt.Accept(in)
	} else {
		in.loopCtrl.doContinue = false
	}

	return nil
}

func (in *interpreter) resolve(expr IExpr, depth int) {
	d := dist{
		distance: depth,
		set:      true,
	}

	in.locals[expr.GetUUID()] = d
}

func (in *interpreter) executeBlock(statements []IStmt, env Environment) any {
	prev := in.env

	// Make sure env is reset at the end, even if something panics
	defer func() {
		in.env = prev
	}()

	in.env = env

	for _, s := range statements {
		in.execute(s)

		if in.env.returnVal != nil {
			rVal := in.env.returnVal
			in.env.returnVal = nil
			return rVal
		}
	}

	return nil
}

func (in *interpreter) VisitForBlockStmt(stmt *Block) any {
	env := LocalEnv(in.env)
	val := in.executeBlock(stmt.Statements, env)

	in.env.returnVal = val

	return nil
}

func (in *interpreter) VisitForExpressionStmt(stmt *Expression) any {
	in.evaluate(stmt.Expression)

	return nil
}

func (in *interpreter) VisitForFunctionStmt(stmt *Function) any {
	fn := NewFunction(stmt, in.env)
	in.env.Define(stmt.Name.Lexeme(), fn)

	return nil
}

func (in *interpreter) VisitForIfStmt(stmt *If) any {
	if isTruthy(in.evaluate(stmt.condition)) {
		in.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		in.execute(stmt.elseBranch)
	}

	return nil
}

func (in *interpreter) VisitForPrintStmt(stmt *Print) any {
	val := in.evaluate(stmt.Expression)
	fmt.Println(stringify(val))

	return nil
}

func (in *interpreter) VisitForReturnStmt(stmt *Return) any {
	if stmt.Value != nil {
		val := in.evaluate(stmt.Value)
		if val != nil {
			in.env.returnVal = val
		}
	}

	return nil
}

func (in *interpreter) VisitForVarStmt(stmt *Var) any {
	var val any

	if stmt.Initialiser != nil {
		val = in.evaluate(stmt.Initialiser)
	}

	in.env.Define(stmt.Name.Lexeme(), val)

	return nil
}

func (in *interpreter) VisitForWhileStmt(stmt *While) any {
	for isTruthy(in.evaluate(stmt.condition)) {
		if in.loopCtrl.doBreak {
			in.loopCtrl.doBreak = false
			return nil
		}

		in.execute(stmt.body)

	}

	return nil
}

func (in *interpreter) VisitForBreakStmt(stmt *Break) any {
	in.loopCtrl.doBreak = true
	return nil
}

func (in *interpreter) VisitForContinueStmt(stmt *Continue) any {
	in.loopCtrl.doContinue = true
	return nil
}

func (in *interpreter) VisitForAssignExpr(expr *Assign) any {

	value := in.evaluate(expr.Value)

	d := in.locals[expr.GetUUID()]

	if d.set {
		in.env.AssignAt(d.distance, expr.Name, value)
	} else {
		in.globals.Assign(expr.Name, value)
	}

	return value
}

func (in *interpreter) VisitForVariableExpr(expr *Variable) any {
	return in.lookUpVariable(expr.Name, expr)
}

func (in *interpreter) lookUpVariable(name IToken, expr IExpr) any {
	d := in.locals[expr.GetUUID()]

	if d.set {
		return in.env.GetAt(d.distance, name.Lexeme())
	}

	return in.globals.Get(name)
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	if _, ok := obj.(bool); ok {
		return obj.(bool)
	}

	return true
}

func isEqual(l any, r any) bool {
	if l == nil && r == nil {
		return true
	}

	if l == nil {
		return false
	}

	return l == r
}

func (in *interpreter) checkNumberOperand(operator IToken, operand any) bool {
	_, floatCheck := operand.(float64)
	if floatCheck {
		return true
	}

	in.handleError(operator, "Operand must be a number.")
	return false
}

func (in *interpreter) checkNumberOperands(operator IToken, left any, right any) bool {
	_, isFloatLeft := left.(float64)
	_, isFloatRight := right.(float64)

	if isFloatLeft && isFloatRight {
		return true
	}

	in.handleError(operator, "Operands must be numbers.")
	return false
}

func (in *interpreter) concatenateOperands(operator IToken, left any, right any) string {

	isLeft, isRight := in.checkStringOperands(operator, left, right)

	var rStr string
	var lStr string

	if isLeft && isRight {
		rStr = right.(string)
		lStr = left.(string)

	}

	if isLeft {
		rStr = stringify(right)
		lStr = left.(string)
	}

	if isRight {
		lStr = stringify(left)
		rStr = right.(string)
	}

	if !isLeft && !isRight {
		in.handleError(operator, "One operand must be a string.")
	}

	return fmt.Sprintf("%s%s", lStr, rStr)

}

// Returns true if either operand can be asserted to a string
func (in *interpreter) checkStringOperands(operator IToken, left any, right any) (bool, bool) {

	_, isStringLeft := left.(string)
	_, isStringRight := right.(string)

	if isStringLeft && isStringRight {
		return true, true
	}

	if isStringLeft {
		return true, false
	}

	if isStringRight {
		return false, true
	}

	return false, false
}

func (in *interpreter) handleError(operator IToken, msg string) {
	e := runtimeError{
		operator,
		msg,
	}

	// Unwind to recover function in Interpret()
	panic(e)
}

func stringify(obj any) string {
	if obj == nil {
		return "nil"
	}

	if _, isFloat := obj.(float64); isFloat {
		text := fmt.Sprintf("%f", obj)

		if hasDecimals(text) {

			return strings.Split(text, ".")[0]
		}
	}

	return fmt.Sprint(obj)
}

func hasDecimals(text string) bool {
	if strings.Contains(text, ".") {
		rhs := strings.Split(text, ".")[1]
		decimals := strings.Split(rhs, "")

		for _, d := range decimals {
			if d != "0" {
				return true
			}
		}

	}

	return false

}

func lrToFloat(left any, right any) (float64, float64, error) {

	fLeft, okLeft := left.(float64)
	if !okLeft {
		return 0.0, 0.0, errors.New("Left operand not float")
	}

	fRight, okRight := right.(float64)
	if !okRight {
		return 0.0, 0.0, errors.New("Right operand not float")
	}

	return fLeft, fRight, nil
}

func anyToFloat(obj any) (float64, error) {
	f, ok := obj.(float64)
	if !ok {
		return 0.0, errors.New("Operand not float")
	}

	return f, nil
}

func anyToFloats(args ...any) ([]float64, []error) {

	floats := []float64{}
	errors := []error{}

	for _, v := range args {
		f, e := strconv.ParseFloat(v.(string), 64)
		if e != nil {
			errors = append(errors, e)
		} else {
			floats = append(floats, f)
		}
	}

	return floats, errors

}
