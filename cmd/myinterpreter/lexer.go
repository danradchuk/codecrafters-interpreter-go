package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var LexerError = errors.New("Error:")

type TokenType int

func (t TokenType) String() string {
	return [...]string{
		"LEFT_PAREN",
		"RIGHT_PAREN",
		"LEFT_BRACE",
		"RIGHT_BRACE",
		"PLUS",
		"MINUS",
		"STAR",
		"DOT",
		"COMMA",
		"SEMICOLON",
		"EQUAL",
		"BANG",
		"BANG_EQUAL",
		"EQUAL_EQUAL",
		"SLASH",
		"LESS",
		"LESS_EQUAL",
		"GREATER",
		"GREATER_EQUAL",
		"STRING",
		"NUMBER",
		"IDENTIFIER",
		"AND",
		"CLASS",
		"ELSE",
		"FALSE",
		"FOR",
		"FUN",
		"IF",
		"NIL",
		"OR",
		"PRINT",
		"RETURN",
		"SUPER",
		"THIS",
		"TRUE",
		"VAR",
		"WHILE",
		"EOF",
	}[t]
}

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	PLUS
	MINUS
	STAR
	DOT
	COMMA
	SEMICOLON
	EQUAL
	BANG
	BANG_EQUAL
	EQUAL_EQUAL
	SLASH
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	STRING
	NUMBER
	IDENTIFIER
	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

func tokenType(c string) TokenType {
	switch c {
	case "(":
		return LEFT_PAREN
	case ")":
		return RIGHT_PAREN
	case "{":
		return LEFT_BRACE
	case "}":
		return RIGHT_BRACE
	case "+":
		return PLUS
	case "-":
		return MINUS
	case "/":
		return SLASH
	case ".":
		return DOT
	case ",":
		return COMMA
	case ";":
		return SEMICOLON
	case "=":
		return EQUAL
	case "!":
		return BANG
	case "!=":
		return BANG_EQUAL
	case "==":
		return EQUAL_EQUAL
	case "*":
		return STAR
	case "<":
		return LESS
	case "<=":
		return LESS_EQUAL
	case ">":
		return GREATER
	case ">=":
		return GREATER_EQUAL
	case "EOF":
		return EOF
	}

	panic("unknown character")
}

var keywordToTokenType = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
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

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

type Lexer struct {
	Input    []byte
	CurrLine int
	CurrPos  int
	ReadPos  int
	Char     rune // Char instead of ch because in Go that usually means channels
}

// Peek observes the next token to consume without advancing to it
func (l *Lexer) Peek() rune {
	if l.ReadPos >= len(l.Input) {
		return 0
	}

	return rune(l.Input[l.ReadPos])
}
func (l *Lexer) skipWhitespaces() {
	for l.Char == ' ' || l.Char == '\t' || l.Char == '\n' || l.Char == '\r' {
		l.CurrPos++
		if l.Char == '\n' || l.Char == '\r' {
			l.CurrLine++
		}
	}
}
func (l *Lexer) ReadChar() {
	if l.ReadPos >= len(l.Input) {
		l.Char = 0
	} else {
		l.Char = rune(l.Input[l.ReadPos])
	}
	l.CurrPos = l.ReadPos
	l.ReadPos++
}
func (l *Lexer) ReadString() string {
	startPos := l.CurrPos + 1
	for {
		l.ReadChar()
		if l.Char == 0 || l.Char == '"' {
			break
		}
	}
	str := string(l.Input[startPos:l.CurrPos])
	return str
}
func (l *Lexer) ReadNumber() string {
	startPos := l.CurrPos
	for unicode.IsDigit(l.Char) {
		l.ReadChar()
	}

	if l.Char == '.' && unicode.IsDigit(l.Peek()) {
		l.ReadChar() // consume '.'

		for unicode.IsDigit(l.Char) {
			l.ReadChar()
		}
	}

	return string(l.Input[startPos:l.CurrPos])
}
func (l *Lexer) ReadIdentifier() string {
	startPos := l.CurrPos
	for isAlphaNumeric(l.Char) {
		l.ReadChar()
	}

	return string(l.Input[startPos:l.CurrPos])
}
func (l *Lexer) NextToken() (*Token, error) {
	var token *Token
	var err error

	l.skipWhitespaces()

	switch l.Char {
	case '(', ')', '{', '}', '+', '-', '*', '.', ',', ';':
		token = &Token{Type: tokenType(string(l.Char)), Lexeme: string(l.Char)}
	case '/':
		if l.Peek() == '/' {
			l.ReadChar()
			for l.Char != '\n' && l.Char != '\r' && l.Char != 0 {
				l.ReadChar()
			}
			l.CurrLine++
		} else {
			token = &Token{Type: tokenType(string(l.Char)), Lexeme: string(l.Char)}
		}
	case '!', '=', '<', '>':
		if l.Peek() == '=' {
			ch := l.Char
			l.ReadChar()
			l := string(ch) + string(l.Char)
			token = &Token{Type: tokenType(l), Lexeme: l}
		} else {
			token = &Token{Type: tokenType(string(l.Char)), Lexeme: string(l.Char)}
		}
	case '"':
		str := l.ReadString()
		if l.Char == 0 { // EOF
			err = fmt.Errorf("[line %d] %w Unterminated string.", l.CurrLine, LexerError)
		} else {
			token = &Token{Type: STRING, Lexeme: `"` + str + `"`, Literal: str}
		}
	case 0:
		return &Token{Type: tokenType("EOF")}, nil
	default:
		if unicode.IsDigit(l.Char) {
			number := l.ReadNumber()

			// 1234. -> 1234.0
			// 200.00 -> 200.0
			// 100.15 -> 100.15 (UNCHANGED)
			trailZeroes := func(s string) string {
				if strings.Contains(s, ".") {
					s = strings.TrimRight(s, "0")
					s = strings.TrimRight(s, ".")
				}

				if !strings.Contains(s, ".") {
					s += ".0"
				}

				return s
			}
			numLiter := trailZeroes(number)

			token = &Token{Type: NUMBER, Lexeme: number, Literal: numLiter}
		} else if isAlphaNumeric(l.Char) {
			ident := l.ReadIdentifier()
			if tok, ok := keywordToTokenType[ident]; ok {
				token = &Token{Type: tok, Lexeme: ident}
			} else {
				token = &Token{Type: IDENTIFIER, Lexeme: ident}
			}
		} else {
			err = fmt.Errorf("[line %d] %w Unexpected character: %c", l.CurrLine, LexerError, l.Char)
		}
	}

	l.ReadChar() // consume next token

	return token, err
}
func (l *Lexer) Tokens() ([]*Token, []error) {
	var tokens []*Token
	var errs []error

	l.ReadChar() // advance the first position

	for {
		t, err := l.NextToken()
		if err != nil {
			errs = append(errs, err)
		}

		if t != nil {
			tokens = append(tokens, t)
		} else {
			continue
		}

		if t.Type == EOF {
			break
		}
	}

	return tokens, errs
}

func isAlphaNumeric(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}
func printTokens(tokens []*Token) {
	handleLiteral := func(s string) string {
		if s == "" {
			return "null"
		} else {
			return s
		}
	}
	for _, tok := range tokens {
		fmt.Printf("%s %s %s\n", tok.Type, tok.Lexeme, handleLiteral(tok.Literal))
	}
}
func processErrors(errs []error) int {
	var isLexerErr = false
	for _, err := range errs {
		if errors.Is(err, LexerError) {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			isLexerErr = true
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}

	if isLexerErr {
		return 65
	}

	return 1
}
