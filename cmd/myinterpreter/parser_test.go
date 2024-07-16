package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	type args struct {
		minBp       int
		fileContent string
	}

	tests := []struct {
		name string
		args args
		want ASTNode
	}{
		{"parseBool", args{0, "true"}, BooleanLiteral{Token{TRUE, "true", ""}, true}},
		{"parseNil", args{0, "nil"}, NilLiteral{}},
		{"parseNumber", args{0, "42.47"}, NumLiteral{Token{NUMBER, "42.47", "42.47"}, 42.47}},
		{"parseString", args{0, "\"hello\""}, StringLiteral{Token{STRING, "\"hello\"", "hello"}, "hello"}},
		{"parseGroupedExpr", args{0, "(\"hello\")"},
			GroupedExpr{
				Token{Type: LEFT_PAREN, Lexeme: "("},
				StringLiteral{
					Token{STRING, "\"hello\"", "hello"},
					"hello"},
			},
		},
		{"parsePrefixExpr", args{0, "!true"},
			PrefixExpr{
				Token{Type: BANG, Lexeme: "!"},
				"!",
				BooleanLiteral{
					Token{TRUE, "true", ""},
					true,
				},
			},
		},
		{"parseInfixExpr", args{0, "1+1*3"},
			InfixExpr{
				Token{PLUS, "+", ""},
				NumLiteral{Token{NUMBER, "1", "1.0"}, 1.},
				"+",
				InfixExpr{
					Token: Token{STAR, "*", ""},
					Left:  NumLiteral{Token{NUMBER, "1", "1.0"}, 1.},
					Op:    "*",
					Right: NumLiteral{Token{NUMBER, "3", "3.0"}, 3.},
				},
			},
		},
		{"parseInfixExpr", args{0, "-(-58 + 68) * (40 * 40) / (72 + 39)"},
			InfixExpr{
				Token{PLUS, "/", ""},
				InfixExpr{
					Token: Token{STAR, "*", ""},
					Left: PrefixExpr{
						Token: Token{MINUS, "-", ""},
						Op:    "-",
						Right: InfixExpr{
							Token: Token{PLUS, "+", ""},
							Left:  NumLiteral{Token{NUMBER, "58", "58.0"}, 58.},
							Op:    "+",
							Right: NumLiteral{Token{NUMBER, "68", "68.0"}, 68.},
						},
					},
					Op: "*",
					Right: InfixExpr{
						Token: Token{STAR, "*", ""},
						Left:  NumLiteral{Token{NUMBER, "40", "40.0"}, 40.},
						Op:    "*",
						Right: NumLiteral{Token{NUMBER, "40", "40.0"}, 40.},
					},
				},
				"/",
				GroupedExpr{
					Token: Token{LEFT_PAREN, "(", ""},
					Value: InfixExpr{
						Token: Token{STAR, "+", ""},
						Left:  NumLiteral{Token{NUMBER, "72", "72.0"}, 72.},
						Op:    "+",
						Right: NumLiteral{Token{NUMBER, "39", "39.0"}, 39.},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, f := prepareTmpFile(t, tt.args.fileContent)
			defer func(name string) { _ = os.Remove(name) }(f.Name())
			defer func(f *os.File) { _ = f.Close() }(f)

			p := NewParser(NewLexer(content))
			//printer := ASTPrinter{w: os.Stdout}
			if got := p.ParseExpr(LOWEST); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			} else {
				//println(got.String())
				//printer.Print(got)
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
