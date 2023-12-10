package jffy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thejff/jffylang/error"
)

type IScanner interface {
	ScanTokens() []IToken
}

type scan struct {
	source string
	tokens []IToken

	keywords map[string]TokenType

	start   int
	current int
	line    int
}

/* type Tokens struct {
} */

func Scanner(source string) IScanner {
	keywords := map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"fun":    FUN,
		"for":    FOR,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}

	var s IScanner = &scan{
		source: source,
		tokens: []IToken{},

		keywords: keywords,

		start:   0,
		current: 0,
		line:    1,
	}

	return s
}

func (s *scan) ScanTokens() []IToken {

	for !s.isAtEnd() {

		// fmt.Printf("Start: %d\n", s.start)
		s.start = s.current
		// fmt.Printf("Current: %d\n", s.current)
		s.scanToken()

	}

	s.tokens = append(s.tokens, Token(EOF, "", nil, s.line))

	return s.tokens

}

func (s *scan) scanToken() {
	c := s.advance()
	var t TokenType

	switch c {

	// Single chars
	case "(":
		s.addToken(LEFT_PAREN, nil)

	case ")":
		s.addToken(RIGHT_PAREN, nil)

	case "{":
		s.addToken(LEFT_BRACE, nil)

	case "}":
		s.addToken(RIGHT_BRACE, nil)

	case ",":
		s.addToken(COMMA, nil)

	case ".":
		s.addToken(DOT, nil)

	case "-":
		s.addToken(MINUS, nil)

	case "+":
		s.addToken(PLUS, nil)

	case ";":
		s.addToken(SEMICOLON, nil)

	case "*":
		s.addToken(STAR, nil)

	// Operators
	case "!":
		t = BANG
		if s.match("=") {
			t = BANG_EQUAL
		}
		s.addToken(t, nil)

	case "=":
		t = EQUAL
		if s.match("=") {
			t = EQUAL_EQUAL
		}
		s.addToken(t, nil)

	case "<":
		t = LESS
		if s.match("=") {
			t = LESS_EQUAL
		}
		s.addToken(t, nil)

	case ">":
		t = GREATER
		if s.match("=") {
			t = GREATER_EQUAL
		}
		s.addToken(t, nil)

	// Division
	case "/":
		if s.match("/") {
			for s.peek() != string('\n') && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}

	// Whitespace, newline etc.
	case " ":
	case string('\r'):
	case string('\t'):
		// Ignore these
		break

	case string('\n'):
		s.line++

	// String literals
	case string('"'):
		s.string()

	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {

			s.identifier()

		} else {
			error.Error(s.line, fmt.Sprintf("Unexpected character: \"%s\"", c))
		}
	}

}

func (s *scan) identifier() {

	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]

	// TODO: smaller algo?
	text = strings.TrimSpace(text)

	tType := s.keywords[text]
	if tType == 0 {
		tType = IDENTIFIER
	}

	s.addToken(tType, nil)
}

func (s *scan) number() {

	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == "." && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	numStr := s.source[s.start:s.current]
	decimal, err := strconv.ParseFloat(numStr, 32)

	if err != nil {
		fmt.Println("Error converting to number")
		fmt.Println(err)
	}

	s.addToken(NUMBER, decimal)
}

func (s *scan) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scan) advance() string {
	c := string(s.source[s.current])
	s.current++
	return c
}

func (s *scan) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}

	if string(s.source[s.current]) != expected {
		return false
	}

	s.current++

	return true
}

func (s *scan) peek() string {
	if s.isAtEnd() {
		return ""
	}

	return string(s.source[s.current])
}

func (s *scan) peekNext() string {
	if s.current+1 >= len(s.source) {
		return ""
	}

	return string(s.source[s.current+1])
}

func (s *scan) addToken(tType TokenType, literal any) {
	text := s.source[s.start:s.current]

	token := Token(tType, text, literal, s.line)
	s.tokens = append(s.tokens, token)
}

func (s *scan) string() {
	for s.peek() != string('"') && !s.isAtEnd() {
		if s.peek() == string('\n') {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		error.Error(s.line, "Unterminated string.")
		return
	}

	// Catch the closing "
	s.advance()

	value := s.source[s.start+1 : s.current]

	s.addToken(STRING, value)
}

func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || c == "_"
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}
