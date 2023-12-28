package jffy

type ICallable interface {
	Arity() int
	Call(i *interpreter, args []any) any
	ToString() string
}
