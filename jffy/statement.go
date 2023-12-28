package jffy

type Stmt struct {}

type IStmt interface {
  Accept(StmtVisitor) any
}

type StmtVisitor interface {
  VisitForBlockStmt(*Block) any
  VisitForExpressionStmt(*Expression) any
  VisitForFunctionStmt(*Function) any
  VisitForIfStmt(*If) any
  VisitForWhileStmt(*While) any
  VisitForBreakStmt(*Break) any
  VisitForContinueStmt(*Continue) any
  VisitForPrintStmt(*Print) any
  VisitForReturnStmt(*Return) any
  VisitForVarStmt(*Var) any
}

type Block struct {
  Statements []IStmt
}

func (b *Block) Accept(param StmtVisitor) any {
  return param.VisitForBlockStmt(b)
}

type Expression struct {
  Expression IExpr
}

func (e *Expression) Accept(param StmtVisitor) any {
  return param.VisitForExpressionStmt(e)
}

type Function struct {
  Name IToken
  Params []IToken
  Body []IStmt
}

func (f *Function) Accept(param StmtVisitor) any {
  return param.VisitForFunctionStmt(f)
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

type Print struct {
  Expression IExpr
}

func (p *Print) Accept(param StmtVisitor) any {
  return param.VisitForPrintStmt(p)
}

type Return struct {
  Keyword IToken
  Value IExpr
}

func (r *Return) Accept(param StmtVisitor) any {
  return param.VisitForReturnStmt(r)
}

type Var struct {
  Name IToken
  Initialiser IExpr
}

func (v *Var) Accept(param StmtVisitor) any {
  return param.VisitForVarStmt(v)
}


