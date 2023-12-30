package jffy

type Expr struct {}

type IExpr interface {
  Accept(ExprVisitor) any
  GetUUID() string
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
  uuid string
}

func (a *Assign) Accept(param ExprVisitor) any {
  return param.VisitForAssignExpr(a)
}

func (a *Assign) GetUUID() string {
  return a.uuid
}

type Binary struct {
  Left IExpr
  Operator IToken
  Right IExpr
  uuid string
}

func (b *Binary) Accept(param ExprVisitor) any {
  return param.VisitForBinaryExpr(b)
}

func (b *Binary) GetUUID() string {
  return b.uuid
}

type Call struct {
  Callee IExpr
  Paren IToken
  Arguments []IExpr
  uuid string
}

func (c *Call) Accept(param ExprVisitor) any {
  return param.VisitForCallExpr(c)
}

func (c *Call) GetUUID() string {
  return c.uuid
}

type Lambda struct {
  Paren IToken
  Params []IToken
  Body []IStmt
  uuid string
}

func (l *Lambda) Accept(param ExprVisitor) any {
  return param.VisitForLambdaExpr(l)
}

func (l *Lambda) GetUUID() string {
  return l.uuid
}

type Grouping struct {
  Expression IExpr
  uuid string
}

func (g *Grouping) Accept(param ExprVisitor) any {
  return param.VisitForGroupingExpr(g)
}

func (g *Grouping) GetUUID() string {
  return g.uuid
}

type Literal struct {
  Value any
  uuid string
}

func (l *Literal) Accept(param ExprVisitor) any {
  return param.VisitForLiteralExpr(l)
}

func (l *Literal) GetUUID() string {
  return l.uuid
}

type Logical struct {
  Left IExpr
  Operator IToken
  Right IExpr
  uuid string
}

func (l *Logical) Accept(param ExprVisitor) any {
  return param.VisitForLogicalExpr(l)
}

func (l *Logical) GetUUID() string {
  return l.uuid
}

type Variable struct {
  Name IToken
  uuid string
}

func (v *Variable) Accept(param ExprVisitor) any {
  return param.VisitForVariableExpr(v)
}

func (v *Variable) GetUUID() string {
  return v.uuid
}

type Unary struct {
  Operator IToken
  Right IExpr
  uuid string
}

func (u *Unary) Accept(param ExprVisitor) any {
  return param.VisitForUnaryExpr(u)
}

func (u *Unary) GetUUID() string {
  return u.uuid
}


