package jffy

type Expr struct {}

type IExpr interface {
  accept(ExprVisitor) any
}

type ExprVisitor interface {
  visitForBinaryExpr(*Binary) any
  visitForGroupingExpr(*Grouping) any
  visitForLiteralExpr(*Literal) any
  visitForUnaryExpr(*Unary) any
}

type Binary struct {
  Left IExpr
  Operator IToken
  Right IExpr
}

func (b *Binary) accept(v ExprVisitor) any {
  return v.visitForBinaryExpr(b)
}

type Grouping struct {
  Expression IExpr
}

func (g *Grouping) accept(v ExprVisitor) any {
  return v.visitForGroupingExpr(g)
}

type Literal struct {
  Value any
}

func (l *Literal) accept(v ExprVisitor) any {
  return v.visitForLiteralExpr(l)
}

type Unary struct {
  Operator IToken
  Right IExpr
}

func (u *Unary) accept(v ExprVisitor) any {
  return v.visitForUnaryExpr(u)
}


