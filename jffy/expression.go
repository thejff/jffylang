package jffy

type Expr struct{}

type IExpr interface {
	Accept(ExprVisitor) any
}

type ExprVisitor interface {
	VisitForBinaryExpr(*Binary) any
	VisitForGroupingExpr(*Grouping) any
	VisitForLiteralExpr(*Literal) any
	VisitForUnaryExpr(*Unary) any
}

type Binary struct {
	Left     IExpr
	Operator IToken
	Right    IExpr
}

func (b *Binary) Accept(v ExprVisitor) any {
	return v.VisitForBinaryExpr(b)
}

type Grouping struct {
	Expression IExpr
}

func (g *Grouping) Accept(v ExprVisitor) any {
	return v.VisitForGroupingExpr(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v ExprVisitor) any {
	return v.VisitForLiteralExpr(l)
}

type Unary struct {
	Operator IToken
	Right    IExpr
}

func (u *Unary) Accept(v ExprVisitor) any {
	return v.VisitForUnaryExpr(u)
}
