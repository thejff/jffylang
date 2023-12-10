package jffy

type Expr struct {
}

type Binary struct {
	left     Expr
	operator IToken
	right    Expr
}
