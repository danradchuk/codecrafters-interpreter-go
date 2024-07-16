package main

import (
	"io"
	"strings"
)

type ASTPrinter struct {
	w   io.Writer
	err error
}

func (v *ASTPrinter) Print(n ASTNode) {
	if n != nil {
		n.Accept(v)
	}

	if v.err != nil {
		panic(v.err)
	}
}

func (v *ASTPrinter) VisitBoolean(n BooleanLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitNil(n NilLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitNum(n NumLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitString(n StringLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitGroupedExpr(n GroupedExpr) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitPrefixExpr(n PrefixExpr) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitInfixExpr(n InfixExpr) {
	v.write(n.String())
}

func (v *ASTPrinter) write(s string) {
	if v.err != nil {
		return
	}
	_, v.err = v.w.Write([]byte(s))
}

func trailZeroes(s string) string {
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	if !strings.Contains(s, ".") {
		s += ".0"
	}

	return s
}
