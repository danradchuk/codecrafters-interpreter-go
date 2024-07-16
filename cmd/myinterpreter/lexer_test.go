package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLexer_Tokens(t *testing.T) {
	type args struct {
		fileContent string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{"scanOps", args{"<"}, []Token{{Type: LESS, Lexeme: "<"}, {Type: EOF}}},
		{"scanOps", args{"2+2=4"}, []Token{
			{Type: NUMBER, Lexeme: "2", Literal: "2.0"},
			{Type: PLUS, Lexeme: "+"},
			{Type: NUMBER, Lexeme: "2", Literal: "2.0"},
			{Type: EQUAL, Lexeme: "="},
			{Type: NUMBER, Lexeme: "4", Literal: "4.0"},
			{Type: EOF},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, f := prepareTmpFile(t, tt.args.fileContent)
			defer func(name string) { _ = os.Remove(name) }(f.Name())
			defer func(f *os.File) { _ = f.Close() }(f)

			l := NewLexer(content)
			if got := l.Tokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokens() = %v, want %v", got, tt.want)
			} else {
				//PrintTokens(got)
			}
		})
	}
}
