package main

import (
	"errors"
	"fmt"
	"os"
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
		"COMMENT",
		"ERROR",
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
	COMMENT
	ERROR
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
	Line    int
}

type Lexer struct {
	input    []byte
	errs     []error
	currLine int
	currPos  int
	readPos  int
	char     rune
}

func NewLexer(input []byte) *Lexer {
	l := &Lexer{
		input:    input,
		currLine: 1,
	}
	l.readChar() // advance cursor to the first position

	return l
}

func (l *Lexer) NextToken() Token {
	var token Token

	l.skipWhitespaces()

	switch l.char {
	case '(', ')', '{', '}', '+', '-', '*', '.', ',', ';':
		token = Token{Type: tokenType(string(l.char)), Lexeme: string(l.char), Line: l.currLine}
	case '/':
		if l.peek() == '/' {
			for l.char != '\n' && l.char != '\r' && l.char != 0 {
				l.readChar()
			}
			if l.char == '\n' || l.char == '\r' {
				l.currLine++
			}
			token = Token{Type: COMMENT}
		} else {
			token = Token{Type: tokenType(string(l.char)), Lexeme: string(l.char), Line: l.currLine}
		}
	case '!', '=', '<', '>':
		if l.peek() == '=' {
			ch := l.char
			l.readChar()
			lex := string(ch) + string(l.char)
			token = Token{Type: tokenType(lex), Lexeme: lex, Line: l.currLine}
		} else {
			token = Token{Type: tokenType(string(l.char)), Lexeme: string(l.char), Line: l.currLine}
		}
	case '"':
		str := l.readString()
		if l.char == 0 { // EOF
			l.errs = append(l.errs, fmt.Errorf("[line %d] %w Unterminated string.", l.currLine, LexerError))
			token = Token{Type: ERROR, Lexeme: string(l.char)}
		} else {
			token = Token{Type: STRING, Lexeme: `"` + str + `"`, Literal: str, Line: l.currLine}
		}
	case 0:
		token = Token{Type: tokenType("EOF")}
	default:
		if unicode.IsDigit(l.char) {
			number := l.readNumber()

			// 1234. -> 1234.0
			// 200.00 -> 200.0
			// 100.15 -> 100.15 (UNCHANGED)

			token = Token{Type: NUMBER, Lexeme: number, Literal: trailZeroes(number), Line: l.currLine}
			return token
		} else if isAlphaNumeric(l.char) {
			ident := l.readIdentifier()
			if tok, ok := keywordToTokenType[ident]; ok {
				token = Token{Type: tok, Lexeme: ident, Line: l.currLine}
			} else {
				token = Token{Type: IDENTIFIER, Lexeme: ident, Line: l.currLine}
			}
			return token
		} else {
			l.errs = append(l.errs, fmt.Errorf("[line %d] %w Unexpected character: %c", l.currLine, LexerError, l.char))
			token = Token{Type: ERROR, Lexeme: string(l.char)}
		}
	}

	l.readChar() // consume next token

	return token
}
func (l *Lexer) Tokens() []Token {
	var tokens []Token
	for t := l.NextToken(); ; t = l.NextToken() {
		if t.Type != COMMENT && t.Type != ERROR {
			tokens = append(tokens, t)
		}

		if t.Type == EOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) peek() rune {
	if l.readPos >= len(l.input) {
		return 0
	}

	return rune(l.input[l.readPos])
}
func (l *Lexer) skipWhitespaces() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.char == '\n' || l.char == '\r' {
			l.currLine++
		}
		l.readChar()
	}
}
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = rune(l.input[l.readPos])
	}
	l.currPos = l.readPos
	l.readPos++
}
func (l *Lexer) readString() string {
	startPos := l.currPos + 1
	for {
		l.readChar()
		if l.char == 0 || l.char == '"' {
			break
		}
	}
	str := string(l.input[startPos:l.currPos])

	return str
}
func (l *Lexer) readNumber() string {
	startPos := l.currPos
	for unicode.IsDigit(l.char) {
		l.readChar()
	}

	if l.char == '.' && unicode.IsDigit(l.peek()) {
		l.readChar() // consume '.'

		for unicode.IsDigit(l.char) {
			l.readChar()
		}
	}

	return string(l.input[startPos:l.currPos])
}
func (l *Lexer) readIdentifier() string {
	startPos := l.currPos
	for isAlphaNumeric(l.char) {
		l.readChar()
	}

	return string(l.input[startPos:l.currPos])
}

func isAlphaNumeric(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}
func PrintTokens(tokens []Token) {
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
func CheckLexerErrors(errs []error) int {
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
