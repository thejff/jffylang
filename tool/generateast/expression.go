package jffy

type Expr struct {}

type IExpr interface {
  accept(visitor)
}

type visitor interface {
  visitForBinaryExpr(*Binary)
  visitForGroupingExpr(*Grouping)
  visitForLiteralExpr(*Literal)
  visitForUnaryExpr(*Unary)
}

type Binary struct {
  left Expr
  operator IToken
  right Expr
}

func (b *Binary) accept(v visitor) {
  v.visitForBinary(b)
}

type Grouping struct {
  expression Expr
}

func (g *Grouping) accept(v visitor) {
  v.visitForGrouping(g)
}

type Literal struct {
  value any
}

func (l *Literal) accept(v visitor) {
  v.visitForLiteral(l)
}

type Unary struct {
  operator IToken
  right Expr
}

func (u *Unary) accept(v visitor) {
  v.visitForUnary(u)
}


