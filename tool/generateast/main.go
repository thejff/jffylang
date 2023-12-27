package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
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

	defineExpressions(outDir)
	defineStatements(outDir)

}

func defineExpressions(outDir string) {
	t := []Types{
		{
			name:   "Assign",
			fields: []string{"Name IToken", "Value IExpr"},
		},
		{
			name:   "Binary",
			fields: []string{"Left IExpr", "Operator IToken", "Right IExpr"},
		},
		{
			name:   "Grouping",
			fields: []string{"Expression IExpr"},
		},
		{
			name:   "Literal",
			fields: []string{"Value any"},
		},
		{
			name:   "Logical",
			fields: []string{"Left IExpr", "Operator IToken", "Right IExpr"},
		},
		{
			name:   "Variable",
			fields: []string{"Name IToken"},
		},
		{
			name:   "Unary",
			fields: []string{"Operator IToken", "Right IExpr"},
		},
	}

	if err := defineAst(outDir, "expression", "Expr", t); err != nil {
		log.Println("Error writing expressions AST")
		log.Fatalln(err)
	}
}

func defineStatements(outDir string) {
	t := []Types{
		{
			name:   "Block",
			fields: []string{"Statements []IStmt"},
		},
		{
			name:   "StmtExpression",
			fields: []string{"Expression IExpr"},
		},
		{
			name:   "If",
			fields: []string{"condition IExpr", "thenBranch IStmt", "elseBranch IStmt"},
		},
		{
			name:   "While",
			fields: []string{"condition IExpr", "body IStmt"},
		},
		{
			name:   "Break",
			fields: []string{},
		},
		{
			name:   "Continue",
			fields: []string{},
		},
		{
			name:   "Var",
			fields: []string{"Name IToken", "Initialiser IExpr"},
		},
		{
			name:   "StmtPrint",
			fields: []string{"Expression IExpr"},
		},
	}

	if err := defineAst(outDir, "statement", "Stmt", t); err != nil {
		log.Println("Error writing statement AST")
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
	lw.writeLine(fmt.Sprintf("type I%s interface {", baseName))
	lw.writeLine(fmt.Sprintf("  Accept(%sVisitor) any", baseName))
	lw.writeLine("}")
	lw.writeLine("")

	defineVisitor(lw, baseName, types)

	for _, t := range types {
		defineStruct(lw, t.name, t.fields, baseName)
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

func defineStruct(lw lineWriter, name string, fields []string, baseName string) {

	lw.writeLine(fmt.Sprintf("type %s struct {", name))

	for _, f := range fields {
		lw.writeLine(
			fmt.Sprintf("  %s", f),
		)
	}

	lw.writeLine("}")
	lw.writeLine("")

	defineBaseFunc(lw, name, baseName)
}

func defineBaseFunc(lw lineWriter, structName string, baseName string) {
	firstChar := strings.ToLower(string(structName[0]))

	lw.writeLine(
		fmt.Sprintf(
			"func (%s *%s) Accept(param %sVisitor) any {",
			firstChar,
			structName,
			baseName,
		),
	)

	lw.writeLine(
		fmt.Sprintf(
			"  return param.VisitFor%s%s(%s)",
			structName,
			baseName,
			firstChar,
		),
	)

	lw.writeLine("}")
	lw.writeLine("")
}

func defineVisitor(lw lineWriter, baseName string, types []Types) {
	lw.writeLine(fmt.Sprintf("type %sVisitor interface {", baseName))

	for _, t := range types {
		lw.writeLine(
			fmt.Sprintf(
				"  VisitFor%s%s(*%s) any",
				t.name,
				baseName,
				t.name,
			),
		)
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
