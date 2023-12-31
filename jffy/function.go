package jffy

import "fmt"

type AnonDeclaration struct {
	Params []IToken
	Body   []IStmt
}

type Declaration struct {
	Name   IToken
	Params []IToken
	Body   []IStmt
}

type AnonJFFFunction struct {
	declaration AnonDeclaration
	closure     Environment
}

type JFFFunction struct {
	declaration Declaration
	closure     Environment
}

func NewFunction(f *Function, closure Environment) ICallable {
	d := Declaration{
		Name:   f.Name,
		Params: f.Params,
		Body:   f.Body,
	}

	return &JFFFunction{
		d,
		closure,
	}
}

func (f *JFFFunction) Call(i *interpreter, args []any) any {
	env := LocalEnv(f.closure)

	for i := 0; i < len(f.declaration.Params); i++ {
		env.Define(args[i])
	}

	return i.executeBlock(f.declaration.Body, env)
}

func (f *JFFFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *JFFFunction) ToString() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme())
}

func NewAnonymousFunction(l *Lambda, closure Environment) ICallable {
	d := AnonDeclaration{
		l.Params,
		l.Body,
	}
	return &AnonJFFFunction{
		d,
		closure,
	}
}

func (f *AnonJFFFunction) Call(i *interpreter, args []any) any {
	env := LocalEnv(f.closure)

	for i := 0; i < len(f.declaration.Params); i++ {
		env.Define(args[i])
	}

	return i.executeBlock(f.declaration.Body, env)
}

func (f *AnonJFFFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *AnonJFFFunction) ToString() string {
	return fmt.Sprintf("<fn anonymous>")
}
