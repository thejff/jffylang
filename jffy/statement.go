package jffy

type Stmt struct {}

type IStmt interface {
  Accept(StmtVisitor) any
}

type StmtVisitor interface {
  VisitForBlockStmt(*Block) any
  VisitForStmtExpressionStmt(*StmtExpression) any
  VisitForIfStmt(*If) any
  VisitForWhileStmt(*While) any
  VisitForBreakStmt(*Break) any
  VisitForContinueStmt(*Continue) any
  VisitForVarStmt(*Var) any
  VisitForStmtPrintStmt(*StmtPrint) any
}

type Block struct {
  Statements []IStmt
}

func (b *Block) Accept(param StmtVisitor) any {
  return param.VisitForBlockStmt(b)
}

type StmtExpression struct {
  Expression IExpr
}

func (s *StmtExpression) Accept(param StmtVisitor) any {
  return param.VisitForStmtExpressionStmt(s)
}

type If struct {
  condition IExpr
  thenBranch IStmt
  elseBranch IStmt
}

func (i *If) Accept(param StmtVisitor) any {
  return param.VisitForIfStmt(i)
}

type While struct {
  condition IExpr
  body IStmt
}

func (w *While) Accept(param StmtVisitor) any {
  return param.VisitForWhileStmt(w)
}

type Break struct {
}

func (b *Break) Accept(param StmtVisitor) any {
  return param.VisitForBreakStmt(b)
}

type Continue struct {
}

func (c *Continue) Accept(param StmtVisitor) any {
  return param.VisitForContinueStmt(c)
}

type Var struct {
  Name IToken
  Initialiser IExpr
}

func (v *Var) Accept(param StmtVisitor) any {
  return param.VisitForVarStmt(v)
}

type StmtPrint struct {
  Expression IExpr
}

func (s *StmtPrint) Accept(param StmtVisitor) any {
  return param.VisitForStmtPrintStmt(s)
}


