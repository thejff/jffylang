package jffy

type Stmt struct {}

type IStmt interface {
  Accept(StmtVisitor) any
}

type StmtVisitor interface {
  VisitForBlockStmt(*Block) any
  VisitForStmtExpressionStmt(*StmtExpression) any
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


