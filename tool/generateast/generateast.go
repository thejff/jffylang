package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
)

type Types struct {
	name   string
	fields []string
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Usage: generateast <output dir>")
		os.Exit(64)
	}

	outDir := args[0]

	t := []Types{
		{
			name:   "Binary",
			fields: []string{"left Expr", "operator IToken", "right Expr"},
		},
		{
			name:   "Grouping",
			fields: []string{"expression Expr"},
		},
		{
			name:   "Literal",
			fields: []string{"value any"},
		},
		{
			name:   "Unary",
			fields: []string{"operator IToken", "right Expr"},
		},
	}

	if err := defineAst(outDir, "expression", "Expr", t); err != nil {
		log.Fatalln(err)
	}
}

func defineAst(outDir string, fileName string, baseName string, types []Types) error {

	file := fmt.Sprintf("%s.go", fileName)
	path := path.Join(outDir, file)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	lw := lineWriter{
		f,
		w,
		[]error{},
		0,
	}

	lw.writeLine("package jffy")
	lw.writeLine("")
	lw.writeLine(fmt.Sprintf("type %s struct {}", baseName))
	lw.writeLine("")

	for _, t := range types {
		defineStruct(lw, t.name, t.fields)
	}

	lw.writeLine("")

	lw.w.Flush()

	if len(lw.errors) > 0 {
		log.Println("Errors occured writing lines to the file")
		for _, e := range lw.errors {
			log.Println(e)
		}
	}

	log.Printf("Total bytes written: %d\n", lw.byteCount)
	return nil
}

func defineStruct(lw lineWriter, name string, fields []string) {

	lw.writeLine(fmt.Sprintf("type %s struct {", name))

	for _, f := range fields {
		lw.writeLine(f)
	}

	lw.writeLine("}")
	lw.writeLine("")
}

type lineWriter struct {
	f         *os.File
	w         *bufio.Writer
	errors    []error
	byteCount int
}

func (l *lineWriter) writeLine(line string) {
	bytesWritten, err := l.w.WriteString(fmt.Sprintf("%s\n", line))
	if err != nil {
		l.errors = append(l.errors, err)
	}

	l.byteCount += bytesWritten
}
