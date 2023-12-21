package jffy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type runtimeError struct {
	operator IToken
	msg      string
}

type interpreter struct {
	jffy Jffy
}

/* func Interpreter(jffy Jffy) ExprVisitor {
	var i ExprVisitor = &interpreter{
		jffy,
	}

	return i
} */

func (in *interpreter) Interpret(expr IExpr, jffy Jffy) {
	in.jffy = jffy

	// Catch panics
	defer func() {
		if r := recover(); r != nil {
			resp, ok := r.(runtimeError)
			if ok {
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

	value := in.evaluate(expr)

	fmt.Printf("%s\n", stringify(value))
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

func (in *interpreter) VisitForGroupingExpr(g *Grouping) any {
	return in.evaluate(g.Expression)
}

func (in *interpreter) VisitForLiteralExpr(l *Literal) any {
	return l.Value
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

	// in.handleError(operator, "Operands must be strings.")
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
