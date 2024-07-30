package eval

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/ast"
)

type Object interface {
	String() string
}

type BooleanObject struct {
	Value bool
}

func (o BooleanObject) String() string {
	return fmt.Sprintf("%t", o.Value)
}

type NilObject struct {
}

func (o NilObject) String() string {
	return "nil"
}

type NumObject struct {
	Value float64
}

func (o NumObject) String() string {
	return trailZeroes(fmt.Sprintf("%f", o.Value))
}

type StrObject struct {
	Value string
}

func (o StrObject) String() string {
	return fmt.Sprintf("%s", o.Value)
}

type GroupedObject struct {
	Value Object
}

func (o GroupedObject) String() string {
	return fmt.Sprintf("%s", o.Value.String())
}

type Evaluator struct {
}

func (e Evaluator) VisitBoolean(n ast.BooleanLiteral) interface{} {
	return &BooleanObject{Value: n.Value}
}
func (e Evaluator) VisitNil(_ ast.NilLiteral) interface{} {
	return &NilObject{}
}
func (e Evaluator) VisitNum(node ast.NumLiteral) interface{} {
	return &NumObject{Value: node.Value}
}
func (e Evaluator) VisitString(node ast.StringLiteral) interface{} {
	return &StrObject{Value: node.Value}
}
func (e Evaluator) VisitGroupedExpr(node ast.GroupedExpr) interface{} {
	expr := node.Value.Accept(e)
	return &GroupedObject{Value: expr.(Object)}
}
func (e Evaluator) VisitPrefixExpr(node ast.PrefixExpr) interface{} {
	expr := node.Right.Accept(e)
	if expr == nil {
		panic("can't evaluate prefix expression")
	}

	switch node.Op {
	case "-":
		if expr, ok := expr.(*NumObject); ok {
			return &NumObject{Value: -expr.Value}
		}
	case "!":
		if _, ok := expr.(*NilObject); ok {
			return &BooleanObject{Value: true}
		}

		if _, ok := expr.(*NumObject); ok {
			return &BooleanObject{Value: false}
		}

		if expr, ok := expr.(*BooleanObject); ok {
			switch expr.Value {
			case true:
				return &BooleanObject{Value: false}
			case false:
				return &BooleanObject{Value: true}
			}
		}
	}

	return nil
}
func (e Evaluator) VisitInfixExpr(node ast.InfixExpr) interface{} {
	return nil
}

func (e Evaluator) Eval(tree ast.Node) Object {
	expr := tree.Accept(e)
	if t, ok := expr.(Object); ok == true {
		return t
	}

	return nil
}

func trailZeroes(s string) string {
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	return s
}
