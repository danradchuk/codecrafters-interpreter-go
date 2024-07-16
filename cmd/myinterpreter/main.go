package main

import (
	"fmt"
	"os"
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
		l := NewLexer(fileContents)
		tokens := l.Tokens()
		PrintTokens(tokens)
		if len(l.errs) > 0 {
			code := CheckLexerErrors(l.errs)
			os.Exit(code)
		}
	} else {
		// Tokenize
		l := NewLexer(fileContents)

		// Parse
		p := NewParser(l)
		ast := p.ParseExpr(LOWEST)
		if len(p.errs) > 0 {
			code := CheckParserErrors(p.errs)
			os.Exit(code)
		}

		// Print
		fmt.Println(ast.String())
		//printer := ASTPrinter{w: os.Stdout}
		//printer.Print(ast)
	}
}
