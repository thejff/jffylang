package jffy

import "fmt"

type IToken interface {
	String() string
	Lexeme() string
	Type() TokenType
	Literal() any
	Line() int
}

type tok struct {
	tType   TokenType
	lexeme  string
	literal any
	line    int
}

func Token(tType TokenType, lexeme string, literal any, line int) IToken {
	var t IToken = &tok{
		tType,
		lexeme,
		literal,
		line,
	}

	return t
}

func (t *tok) String() string {
	// return fmt.Sprintf("type %d lexeme %s literal %v\n", t.tType, t.lexeme, t.literal)
	return fmt.Sprintf("Type %d \nLexeme %s \nLiteral %v\n", t.tType, t.lexeme, t.literal)
}

func (t *tok) Lexeme() string {
	return t.lexeme
}

func (t *tok) Type() TokenType {
	return t.tType
}

func (t *tok) Literal() any {
	return t.literal
}

func (t *tok) Line() int {
	return t.line
}
