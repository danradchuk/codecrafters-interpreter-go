package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
)

func main() {
	if len(os.Args) < 3 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "tokenize" && command != "parse" {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if command == "tokenize" {
		l := lexer.NewLexer(fileContents)
		tokens := l.Tokens()
		lexer.PrintTokens(tokens)
		if len(l.Errors) > 0 {
			code := lexer.CheckErrors(l.Errors)
			os.Exit(code)
		}
	} else if command == "parse" {
		// Tokenize
		l := lexer.NewLexer(fileContents)

		// Parse
		p := parser.NewParser(l)
		ast := p.ParseExpr(parser.LOWEST)
		if len(p.Errors) > 0 {
			code := parser.CheckErrors(p.Errors)
			os.Exit(code)
		}

		// Print
		fmt.Println(ast.String())
	}
}
