package parser

import (
	"os"
	"reflect"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/ast"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
)

func TestParser_Parse(t *testing.T) {
	type args struct {
		minBp       int
		fileContent string
	}

	tests := []struct {
		name string
		args args
		want ast.Node
	}{
		{"parseBool", args{0, "true"},
			ast.BooleanLiteral{
				Token: lexer.Token{Type: lexer.TRUE, Lexeme: "true", Line: 1},
				Value: true,
			},
		},
		{"parseNil", args{0, "nil"}, ast.NilLiteral{}},
		{"parseNumber", args{0, "42.47"},
			ast.NumLiteral{
				Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "42.47", Literal: "42.47", Line: 1},
				Value: 42.47,
			},
		},
		{"parseString", args{0, "\"hello\""},
			ast.StringLiteral{
				Token: lexer.Token{Type: lexer.STRING, Lexeme: "\"hello\"", Literal: "hello", Line: 1},
				Value: "hello",
			},
		},
		{"parseGroupedExpr", args{0, "(\"hello\")"},
			ast.GroupedExpr{
				Token: lexer.Token{Type: lexer.LEFT_PAREN, Lexeme: "(", Line: 1},
				Value: ast.StringLiteral{
					Token: lexer.Token{Type: lexer.STRING, Lexeme: "\"hello\"", Literal: "hello", Line: 1},
					Value: "hello",
				},
			},
		},
		{"parsePrefixExpr", args{0, "!true"},
			ast.PrefixExpr{
				Token: lexer.Token{Type: lexer.BANG, Lexeme: "!", Line: 1},
				Op:    "!",
				Right: ast.BooleanLiteral{
					Token: lexer.Token{Type: lexer.TRUE, Lexeme: "true", Line: 1},
					Value: true,
				},
			},
		},
		{"parseInfixExpr", args{0, "1+1*3"},
			ast.InfixExpr{
				Token: lexer.Token{Type: lexer.PLUS, Lexeme: "+", Line: 1},
				Left:  ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "1", Literal: "1.0", Line: 1}, Value: 1.},
				Op:    "+",
				Right: ast.InfixExpr{
					Token: lexer.Token{Type: lexer.STAR, Lexeme: "*", Line: 1},
					Left:  ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "1", Literal: "1.0", Line: 1}, Value: 1.},
					Op:    "*",
					Right: ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "3", Literal: "3.0", Line: 1}, Value: 3.},
				},
			},
		},
		{"parseInfixExpr", args{0, "-(-58 + 68) * (40 * 40) / (72 + 39)"},
			ast.InfixExpr{
				Token: lexer.Token{Type: lexer.SLASH, Lexeme: "/", Line: 1},
				Left: ast.InfixExpr{
					Token: lexer.Token{Type: lexer.STAR, Lexeme: "*", Line: 1},
					Left: ast.PrefixExpr{
						Token: lexer.Token{Type: lexer.MINUS, Lexeme: "-", Line: 1},
						Op:    "-",
						Right: ast.GroupedExpr{
							Token: lexer.Token{Type: lexer.LEFT_PAREN, Lexeme: "(", Line: 1},
							Value: ast.InfixExpr{
								Token: lexer.Token{Type: lexer.PLUS, Lexeme: "+", Line: 1},
								Left: ast.PrefixExpr{
									Token: lexer.Token{Type: lexer.MINUS, Lexeme: "-", Line: 1},
									Op:    "-",
									Right: ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "58", Literal: "58.0", Line: 1}, Value: 58.},
								},
								Op:    "+",
								Right: ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "68", Literal: "68.0", Line: 1}, Value: 68.},
							},
						},
					},
					Op: "*",
					Right: ast.GroupedExpr{
						Token: lexer.Token{Type: lexer.LEFT_PAREN, Lexeme: "(", Line: 1},
						Value: ast.InfixExpr{
							Token: lexer.Token{Type: lexer.STAR, Lexeme: "*", Line: 1},
							Left:  ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "40", Literal: "40.0", Line: 1}, Value: 40.},
							Op:    "*",
							Right: ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "40", Literal: "40.0", Line: 1}, Value: 40.},
						},
					},
				},
				Op: "/",
				Right: ast.GroupedExpr{
					Token: lexer.Token{Type: lexer.LEFT_PAREN, Lexeme: "(", Line: 1},
					Value: ast.InfixExpr{
						Token: lexer.Token{Type: lexer.PLUS, Lexeme: "+", Line: 1},
						Left:  ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "72", Literal: "72.0", Line: 1}, Value: 72.},
						Op:    "+",
						Right: ast.NumLiteral{Token: lexer.Token{Type: lexer.NUMBER, Lexeme: "39", Literal: "39.0", Line: 1}, Value: 39.},
					},
				},
			},
		},
		{"parseEquality", args{0, "\"foo\" == \"foo\""},
			ast.InfixExpr{
				Token: lexer.Token{Type: lexer.EQUAL_EQUAL, Lexeme: "==", Line: 1},
				Left: ast.StringLiteral{
					Token: lexer.Token{Type: lexer.STRING, Lexeme: "\"foo\"", Literal: "foo", Line: 1},
					Value: "foo",
				},
				Op: "==",
				Right: ast.StringLiteral{
					Token: lexer.Token{Type: lexer.STRING, Lexeme: "\"foo\"", Literal: "foo", Line: 1},
					Value: "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, f := prepareTmpFile(t, tt.args.fileContent)
			defer func(name string) { _ = os.Remove(name) }(f.Name())
			defer func(f *os.File) { _ = f.Close() }(f)

			p := NewParser(lexer.NewLexer(content))
			if got := p.ParseExpr(LOWEST); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func prepareTmpFile(t *testing.T, content string) ([]byte, *os.File) {
	f, err := os.CreateTemp("/tmp", "content")
	if err != nil {
		t.Fatal(err)
	}

	n, err := f.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	if n == 0 {
		t.Fatal("nothing written")
	}

	str, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	return str, f
}
