package main

import (
	"log"
	"os"

	"github.com/thejff/jffylang/jffy"
)

func main() {
	args := os.Args[1:]

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
