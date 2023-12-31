package jffy

import "fmt"

type Environment struct {
	// values    map[string]any
	values    []any
	enclosing *Environment
	returnVal any
}

func GlobalEnv() Environment {
	return Environment{
		// values:    make(map[string]any),
		values:    []any{},
		enclosing: nil,
	}
}

func LocalEnv(enclosing Environment) Environment {
	return Environment{
		values:    []any{},
		enclosing: &enclosing,
	}
}

func (e *Environment) Define(value any) int {
	// e.values[name] = value
	e.values = append(e.values, value)
	return len(e.values) - 1
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e

	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
}

func (e *Environment) GetAt(distance int, index int) any {
	a := e.ancestor(distance)
	return a.values[index]
}

func (e *Environment) AssignAt(distance int, index int, value any) {
	e.ancestor(distance).values[index] = value
}

func (e *Environment) Get(index int, name IToken) any {

	fmt.Printf("This env: \n%v\n\n", e.values)

	if len(e.values) == 0 {
		// Maybe our parent has it?
		if e.enclosing != nil {
			return e.enclosing.Get(index, name)
		}

		// Nobody has it!
		panic(runtimeError{
			operator: name,
			msg:      fmt.Sprintf("Undefined variable \"%s\".", name.Lexeme()),
		})
	}

	val := e.values[index]

	// Found it, send it back
	if val != nil {
		return val
	}

	panic(runtimeError{
		operator: name,
		msg:      fmt.Sprintf("Variable declared but not assigned before use \"%s\".", name.Lexeme()),
	})

}

func (e *Environment) Assign(index int, name IToken, value any) {

	e.values[index] = value

	if e.enclosing != nil {
		e.enclosing.Assign(index, name, value)
		return
	}

	panic(
		runtimeError{
			operator: name,
			msg:      fmt.Sprintf("Undefined variable \"%s\".", name.Lexeme()),
		})
}
