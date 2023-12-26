package jffy

import "fmt"

type Environment struct {
	values    map[string]any
	enclosing *Environment
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

// var a = "ga"; var b = "gb"; var c = "gc"; { var a = "oa"; var b = "ob"; { var a = "ia"; print a; print b; print c; } print a; print b; print c; }

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name IToken) any {
	varName := name.Lexeme()
	val, ok := e.values[varName]

	// Found it, send it back
	if ok {
		return val
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
