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
	return expr.accept(a).(string)
}

func (a *AstPrinter) visitForBinaryExpr(b *Binary) any {
	return a.parenthesize(b.Operator.Lexeme(), b.Left, b.Right)
}

func (a *AstPrinter) visitForGroupingExpr(g *Grouping) any {
	return a.parenthesize("group", g.Expression)
}

func (a *AstPrinter) visitForLiteralExpr(l *Literal) any {
	if l.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", l.Value)
}

func (a *AstPrinter) visitForUnaryExpr(u *Unary) any {
	return a.parenthesize(u.Operator.Lexeme(), u.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...IExpr) string {

	var s strings.Builder

	s.WriteString("(")
	s.WriteString(name)

	for _, e := range exprs {
		s.WriteString(" ")

		child := e.accept(a)
		s.WriteString(child.(string))
	}

	s.WriteString(")")

	return s.String()
}
