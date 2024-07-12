package main

import "fmt"

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
		"EQ",
		"DOUBLE_EQ",
		"DIVISION",
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
	EQ
	DOUBLE_EQ
	DIVISION
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
		return STAR
	case ".":
		return DOT
	case ",":
		return COMMA
	case ";":
		return SEMICOLON
	case "=":
		return EQ
	case "==":
		return DOUBLE_EQ
	case "*":
		return STAR
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

func lexify(input []byte) ([]Token, error) {
	n := len(input)

	var (
		tokens  []Token
		currPos int
	)
	for currPos < n {
		ch := rune(input[currPos]) // TODO proper handling of UTF-8 symbols, for now assume that the input in ASCII encoding

		// skip whitespaces
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			currPos++
			continue
		}

		switch ch {
		case '(', ')', '{', '}', '+', '-', '*', '.', ',', ';':
			tokens = append(tokens, Token{Type: tokenType(string(ch)), Lexeme: string(ch)})
			currPos++
		case '=':
			if peek(input, currPos) == '=' {
				l := string(ch) + string(rune(input[currPos+1]))
				tokens = append(tokens, Token{Type: tokenType(l), Lexeme: l})
				currPos++
			} else {
				tokens = append(tokens, Token{Type: tokenType(string(ch)), Lexeme: string(ch)})
			}
			currPos++
		case '/':
			if peek(input, currPos) == '/' {
				// handle comment
				// read until /n or /r or /r/n
				// advance the position on the count of the skipped symbols
				ahead := handleComment(input, currPos)
				currPos += ahead
			} else {
				tokens = append(tokens, Token{Type: tokenType(string(ch)), Lexeme: string(ch)})
			}
			currPos++
		default:
			currPos++
		}
	}

	tokens = append(tokens, Token{Type: tokenType("EOF")})

	return tokens, nil
}

func peek(input []byte, pos int) rune {
	return rune(input[pos+1])
}

func handleComment(input []byte, pos int) int {
	var ahead int
	for pos < len(input) {
		ch := rune(input[pos])
		if ch == '\n' || ch == '\r' {
			break
		}
		ahead++
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
