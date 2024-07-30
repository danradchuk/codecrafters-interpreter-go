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

func (v *ASTPrinter) VisitBoolean(n ast.BooleanLiteral) interface{} {
	v.write(n.String())
	return nil
}
func (v *ASTPrinter) VisitNil(n ast.NilLiteral) interface{} {
	v.write(n.String())
	return nil
}
func (v *ASTPrinter) VisitNum(n ast.NumLiteral) interface{} {
	v.write(n.String())
	return nil
}
func (v *ASTPrinter) VisitString(n ast.StringLiteral) interface{} {
	v.write(n.String())
	return nil
}
func (v *ASTPrinter) VisitGroupedExpr(n ast.GroupedExpr) interface{} {
	v.write(n.String())
	return nil
}
func (v *ASTPrinter) VisitPrefixExpr(n ast.PrefixExpr) interface{} {
	v.write(n.String())
	return nil
}
func (v *ASTPrinter) VisitInfixExpr(n ast.InfixExpr) interface{} {
	v.write(n.String())
	return nil
}

func (v *ASTPrinter) write(s string) {
	if v.err != nil {
		return
	}
	_, v.err = v.w.Write([]byte(s))
}
