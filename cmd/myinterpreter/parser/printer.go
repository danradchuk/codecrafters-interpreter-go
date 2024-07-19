package parser

import (
	"io"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/ast"
)

type ASTPrinter struct {
	w   io.Writer
	err error
}

func (v *ASTPrinter) Print(n ast.Node) error {
	if n != nil {
		n.Accept(v)
	}

	if v.err != nil {
		return v.err
	}

	return nil
}

func (v *ASTPrinter) VisitBoolean(n ast.BooleanLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitNil(n ast.NilLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitNum(n ast.NumLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitString(n ast.StringLiteral) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitGroupedExpr(n ast.GroupedExpr) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitPrefixExpr(n ast.PrefixExpr) {
	v.write(n.String())
}
func (v *ASTPrinter) VisitInfixExpr(n ast.InfixExpr) {
	v.write(n.String())
}

func (v *ASTPrinter) write(s string) {
	if v.err != nil {
		return
	}
	_, v.err = v.w.Write([]byte(s))
}
