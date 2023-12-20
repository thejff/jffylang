package jffy

type TokenType int

const (
	// Single char tokens
	LEFT_PAREN  TokenType = 0
	RIGHT_PAREN TokenType = 1
	LEFT_BRACE  TokenType = 2
	RIGHT_BRACE TokenType = 3
	COMMA       TokenType = 4
	DOT         TokenType = 5
	MINUS       TokenType = 6
	PLUS        TokenType = 7
	SEMICOLON   TokenType = 8
	SLASH       TokenType = 9
	STAR        TokenType = 10

	// One/Two char tokens
	BANG          TokenType = 11
	BANG_EQUAL    TokenType = 12
	EQUAL         TokenType = 13
	EQUAL_EQUAL   TokenType = 14
	GREATER       TokenType = 15
	GREATER_EQUAL TokenType = 16
	LESS          TokenType = 17
	LESS_EQUAL    TokenType = 18
	DOT_DOT       TokenType = 19

	// Literals
	IDENTIFIER TokenType = 20
	STRING     TokenType = 21
	NUMBER     TokenType = 22

	// Keywords
	AND    TokenType = 23
	CLASS  TokenType = 24
	ELSE   TokenType = 25
	FALSE  TokenType = 26
	FUN    TokenType = 27
	FOR    TokenType = 28
	IF     TokenType = 29
	NIL    TokenType = 30
	OR     TokenType = 31
	PRINT  TokenType = 32
	RETURN TokenType = 33
	SUPER  TokenType = 34
	THIS   TokenType = 35
	TRUE   TokenType = 36
	VAR    TokenType = 37
	WHILE  TokenType = 38

	EOF TokenType = 39
)
