package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

func lexify(input []byte) ([]Token, []error) {
	n := len(input)

	var (
		tokens  []Token
		errs    []error
		currPos int
		line    int
	)

	for currPos < n {
		ch := rune(input[currPos]) // TODO proper handling of UTF-8 symbols, for now assume that the input is in ASCII encoding

		// skip whitespaces
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			currPos++
			if ch == '\n' || ch == '\r' {
				line++
			}
			continue
		}

		switch ch {
		case '(', ')', '{', '}', '+', '-', '*', '.', ',', ';':
			tokens = append(tokens, Token{Type: tokenType(string(ch)), Lexeme: string(ch)})
			currPos++
		case '/':
			if peek(input, currPos+1) == '/' {
				ahead := handleComment(input, currPos)
				currPos += ahead
				line++
			} else {
				tokens = append(tokens, Token{Type: tokenType(string(ch)), Lexeme: string(ch)})
			}
			currPos++
		case '!', '=', '<', '>':
			if peek(input, currPos+1) == '=' {
				l := string(ch) + string(rune(input[currPos+1]))
				tokens = append(tokens, Token{Type: tokenType(l), Lexeme: l})
				currPos++
			} else {
				tokens = append(tokens, Token{Type: tokenType(string(ch)), Lexeme: string(ch)})
			}
			currPos++
		default:
			if ch == '"' {
				currPos++ // consume the opening "
				startPos := currPos
				for currPos < len(input) && input[currPos] != '"' {
					currPos++
				}

				if currPos < len(input) && input[currPos] == '"' {
					str := string(input[startPos:currPos])
					tokens = append(tokens, Token{Type: STRING, Lexeme: `"` + str + `"`, Literal: str})
					currPos++ // consume the closing "
				} else {
					errs = append(errs, fmt.Errorf("[line %d] %w Unterminated string.", line+1, LexerError))
				}
			} else if unicode.IsDigit(ch) {
				startPos := currPos
				for unicode.IsDigit(peek(input, currPos)) {
					currPos++
				}

				if peek(input, currPos) == '.' && unicode.IsDigit(peek(input, currPos+1)) {
					currPos++ // consume '.'

					for unicode.IsDigit(peek(input, currPos)) {
						currPos++
					}
				}

				// format number
				number := string(input[startPos:currPos])
				numLiter := number

				numParts := strings.SplitN(number, ".", 2)
				if len(numParts) != 2 {
					num, _ := strconv.ParseFloat(number, 64)
					numLiter = fmt.Sprintf("%.1f", num)
				}

				tokens = append(tokens, Token{Type: NUMBER, Lexeme: number, Literal: numLiter})
			} else {
				errs = append(errs, fmt.Errorf("[line %d] %w Unexpected character: %c", line+1, LexerError, ch))
				currPos++
			}
		}
	}

	tokens = append(tokens, Token{Type: tokenType("EOF")})

	return tokens, errs
}

func peek(input []byte, pos int) rune {
	if pos >= len(input) {
		return 0
	}

	return rune(input[pos])
}

func handleComment(input []byte, pos int) int {
	var ahead int
	for pos < len(input) {
		ch := rune(input[pos])
		if ch == '\n' || ch == '\r' {
			break
		}
		ahead++
		pos++
	}

	return ahead
}

func printTokens(tokens []Token) {
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
