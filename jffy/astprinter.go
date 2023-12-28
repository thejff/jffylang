package jffy

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() ExprVisitor {
	var p ExprVisitor = &AstPrinter{}
	return p
}

func (a *AstPrinter) Print(expr IExpr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitForBinaryExpr(b *Binary) any {
	return a.parenthesize(b.Operator.Lexeme(), b.Left, b.Right)
}

func (a *AstPrinter) VisitForGroupingExpr(g *Grouping) any {
	return a.parenthesize("group", g.Expression)
}

func (a *AstPrinter) VisitForLiteralExpr(l *Literal) any {
	if l.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", l.Value)
}

func (a *AstPrinter) VisitForVariableExpr(v *Variable) any {
	return nil
}

func (a *AstPrinter) VisitForAssignExpr(v *Assign) any {
	return nil
}

func (a *AstPrinter) VisitForLogicalExpr(v *Logical) any {
	return nil
}

func (a *AstPrinter) VisitForCallExpr(v *Call) any {
	return nil
}

func (a *AstPrinter) VisitForLambdaExpr(l *Lambda) any {
	return nil
}

func (a *AstPrinter) VisitForUnaryExpr(u *Unary) any {
	return a.parenthesize(u.Operator.Lexeme(), u.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...IExpr) string {

	var s strings.Builder

	s.WriteString("(")
	s.WriteString(name)

	for _, e := range exprs {
		s.WriteString(" ")

		child := e.Accept(a)
		s.WriteString(child.(string))
	}

	s.WriteString(")")

	return s.String()
}
