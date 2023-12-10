package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thejff/jffylang/jffy"
)

func main() {
	args := os.Args[1:]

	var uR jffy.IExpr = &jffy.Literal{
		Value: 123,
	}

	var left jffy.IExpr = &jffy.Unary{
		Operator: jffy.Token(jffy.MINUS, "-", nil, 1),
		Right:    uR,
	}

	var op jffy.IToken = jffy.Token(jffy.STAR, "*", nil, 1)

	var right jffy.IExpr = &jffy.Grouping{
		Expression: &jffy.Literal{
			Value: 45.67,
		},
	}

	var e1 jffy.IExpr = &jffy.Binary{
		Left:     left,
		Operator: op,
		Right:    right,
	}

	prn := jffy.NewAstPrinter()
	data := prn.(*jffy.AstPrinter).Print(e1)
	fmt.Println(data)

	jlang := jffy.NewJffy()

	var err error
	if len(args) > 1 {
		jlang.Help()
	} else if len(args) == 1 {
		err = jlang.RunFile(args[0])
	} else {
		err = jlang.RunPrompt()
	}

	if err != nil {
		log.Fatalln(err)
	}
}
