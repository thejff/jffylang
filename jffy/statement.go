package jffy

type Stmt struct {}

type IStmt interface {
  Accept(StmtVisitor) any
  GetUUID() string
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
  uuid string
}

func (b *Block) Accept(param StmtVisitor) any {
  return param.VisitForBlockStmt(b)
}

func (b *Block) GetUUID() string {
  return b.uuid
}

type Expression struct {
  Expression IExpr
  uuid string
}

func (e *Expression) Accept(param StmtVisitor) any {
  return param.VisitForExpressionStmt(e)
}

func (e *Expression) GetUUID() string {
  return e.uuid
}

type Function struct {
  Name IToken
  Params []IToken
  Body []IStmt
  uuid string
}

func (f *Function) Accept(param StmtVisitor) any {
  return param.VisitForFunctionStmt(f)
}

func (f *Function) GetUUID() string {
  return f.uuid
}

type If struct {
  condition IExpr
  thenBranch IStmt
  elseBranch IStmt
  uuid string
}

func (i *If) Accept(param StmtVisitor) any {
  return param.VisitForIfStmt(i)
}

func (i *If) GetUUID() string {
  return i.uuid
}

type While struct {
  condition IExpr
  body IStmt
  uuid string
}

func (w *While) Accept(param StmtVisitor) any {
  return param.VisitForWhileStmt(w)
}

func (w *While) GetUUID() string {
  return w.uuid
}

type Break struct {
  uuid string
}

func (b *Break) Accept(param StmtVisitor) any {
  return param.VisitForBreakStmt(b)
}

func (b *Break) GetUUID() string {
  return b.uuid
}

type Continue struct {
  uuid string
}

func (c *Continue) Accept(param StmtVisitor) any {
  return param.VisitForContinueStmt(c)
}

func (c *Continue) GetUUID() string {
  return c.uuid
}

type Print struct {
  Expression IExpr
  uuid string
}

func (p *Print) Accept(param StmtVisitor) any {
  return param.VisitForPrintStmt(p)
}

func (p *Print) GetUUID() string {
  return p.uuid
}

type Return struct {
  Keyword IToken
  Value IExpr
  uuid string
}

func (r *Return) Accept(param StmtVisitor) any {
  return param.VisitForReturnStmt(r)
}

func (r *Return) GetUUID() string {
  return r.uuid
}

type Var struct {
  Name IToken
  Initialiser IExpr
  uuid string
}

func (v *Var) Accept(param StmtVisitor) any {
  return param.VisitForVarStmt(v)
}

func (v *Var) GetUUID() string {
  return v.uuid
}


