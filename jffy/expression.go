package jffy

type Expr struct {}

type IExpr interface {
  Accept(ExprVisitor) any
}

type ExprVisitor interface {
  VisitForAssignExpr(*Assign) any
  VisitForBinaryExpr(*Binary) any
  VisitForGroupingExpr(*Grouping) any
  VisitForLiteralExpr(*Literal) any
  VisitForVariableExpr(*Variable) any
  VisitForUnaryExpr(*Unary) any
}

type Assign struct {
  Name IToken
  Value IExpr
}

func (a *Assign) Accept(param ExprVisitor) any {
  return param.VisitForAssignExpr(a)
}

type Binary struct {
  Left IExpr
  Operator IToken
  Right IExpr
}

func (b *Binary) Accept(param ExprVisitor) any {
  return param.VisitForBinaryExpr(b)
}

type Grouping struct {
  Expression IExpr
}

func (g *Grouping) Accept(param ExprVisitor) any {
  return param.VisitForGroupingExpr(g)
}

type Literal struct {
  Value any
}

func (l *Literal) Accept(param ExprVisitor) any {
  return param.VisitForLiteralExpr(l)
}

type Variable struct {
  Name IToken
}

func (v *Variable) Accept(param ExprVisitor) any {
  return param.VisitForVariableExpr(v)
}

type Unary struct {
  Operator IToken
  Right IExpr
}

func (u *Unary) Accept(param ExprVisitor) any {
  return param.VisitForUnaryExpr(u)
}


