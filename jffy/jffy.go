package jffy

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Jffy interface {
	Help()
	RunFile(file string) error
	RunPrompt() error
}

type jffy struct {
	hadError bool
}

func NewJffy() Jffy {
	var j Jffy = &jffy{
		hadError: false,
	}

	return j
}

func (j *jffy) Help() {
	fmt.Println("Usage: jffy [script]")
	os.Exit(64)
}

func (j *jffy) RunFile(file string) error {

	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("Error reading script %s\n", file)
	}

	if err := j.run(string(bytes)); err != nil {
		log.Printf("Error running script %s\n", file)
	}

	if j.hadError {
		os.Exit(65)
	}

	return nil

}

func (j *jffy) RunPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("JFFY Interpreter v0.0.1 - http://lang.thejustfor.fun/")

	for {
		fmt.Print("jffy> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %x\n", err)
		}

		if line == "" {
			break
		}

		if err := j.run(line); err != nil {
			fmt.Printf("Error: %x\n", err)
		}

		j.hadError = false
	}

	return nil

}

func (j *jffy) run(source string) error {

	scan := Scanner(source)

	tokens := scan.ScanTokens()

	for _, t := range tokens {
		fmt.Printf("Token: %s\n", t.String())
	}

	return nil
}
