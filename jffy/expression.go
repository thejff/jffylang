package jffy

type Expr struct {}

type IExpr interface {
  Accept(ExprVisitor) any
}

type ExprVisitor interface {
  VisitForAssignExpr(*Assign) any
  VisitForBinaryExpr(*Binary) any
  VisitForCallExpr(*Call) any
  VisitForLambdaExpr(*Lambda) any
  VisitForGroupingExpr(*Grouping) any
  VisitForLiteralExpr(*Literal) any
  VisitForLogicalExpr(*Logical) any
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

type Call struct {
  Callee IExpr
  Paren IToken
  Arguments []IExpr
}

func (c *Call) Accept(param ExprVisitor) any {
  return param.VisitForCallExpr(c)
}

type Lambda struct {
  Paren IToken
  Params []IToken
  Body []IStmt
}

func (l *Lambda) Accept(param ExprVisitor) any {
  return param.VisitForLambdaExpr(l)
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

type Logical struct {
  Left IExpr
  Operator IToken
  Right IExpr
}

func (l *Logical) Accept(param ExprVisitor) any {
  return param.VisitForLogicalExpr(l)
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


