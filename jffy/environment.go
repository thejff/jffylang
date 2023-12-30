package jffy

import "fmt"

type Environment struct {
	values    map[string]any
	enclosing *Environment
	returnVal any
}

func GlobalEnv() Environment {
	return Environment{
		values:    make(map[string]any),
		enclosing: nil,
	}
}

func LocalEnv(enclosing Environment) Environment {
	return Environment{
		values:    make(map[string]any),
		enclosing: &enclosing,
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e

	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
}

func (e *Environment) GetAt(distance int, name string) any {
	a := e.ancestor(distance)
	return a.values[name]
}

func (e *Environment) AssignAt(distance int, name IToken, value any) {
	e.ancestor(distance).values[name.Lexeme()] = value
}

func (e *Environment) Get(name IToken) any {
	varName := name.Lexeme()
	val, ok := e.values[varName]

	// Found it, send it back
	if ok {
		if val != nil {
			return val
		}

		panic(runtimeError{
			operator: name,
			msg:      fmt.Sprintf("Variable declared but not assigned before use \"%s\".", name.Lexeme()),
		})
	}

	// Maybe our parent has it?
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	// Nobody has it!
	panic(runtimeError{
		operator: name,
		msg:      fmt.Sprintf("Undefined variable \"%s\".", name.Lexeme()),
	})
}

func (e *Environment) Assign(name IToken, value any) {

	for k := range e.values {
		if k == name.Lexeme() {
			e.values[name.Lexeme()] = value
			return
		}
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	panic(
		runtimeError{
			operator: name,
			msg:      fmt.Sprintf("Undefined variable \"%s\".", name.Lexeme()),
		})
}
