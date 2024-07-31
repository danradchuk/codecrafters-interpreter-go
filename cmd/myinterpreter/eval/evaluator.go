package eval

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/ast"
)

type Object interface {
	Type() string
	String() string
}

type BooleanObject struct {
	Value bool
}

func (o BooleanObject) Type() string {
	return "BOOLEAN_OBJ"
}
func (o BooleanObject) String() string {
	return fmt.Sprintf("%t", o.Value)
}

type NilObject struct {
}

func (o NilObject) Type() string {
	return "NIL_OBJ"
}
func (o NilObject) String() string {
	return "nil"
}

type NumObject struct {
	Value float64
}

func (o NumObject) Type() string {
	return "NUM_OBJ"
}
func (o NumObject) String() string {
	return trailZeroes(fmt.Sprintf("%f", o.Value))
}

type StrObject struct {
	Value string
}

func (o StrObject) Type() string {
	return "STRING_OBJ"
}
func (o StrObject) String() string {
	return fmt.Sprintf("%s", o.Value)
}

type Evaluator struct {
	Errors []error
}

func (e *Evaluator) Eval(tree ast.Node) Object {
	expr := tree.Accept(e)
	if t, ok := expr.(Object); ok == true {
		return t
	}

	return nil
}

func (e *Evaluator) VisitBoolean(n ast.BooleanLiteral) interface{} {
	return &BooleanObject{Value: n.Value}
}
func (e *Evaluator) VisitNil(_ ast.NilLiteral) interface{} {
	return &NilObject{}
}
func (e *Evaluator) VisitNum(node ast.NumLiteral) interface{} {
	return &NumObject{Value: node.Value}
}
func (e *Evaluator) VisitString(node ast.StringLiteral) interface{} {
	return &StrObject{Value: node.Value}
}
func (e *Evaluator) VisitGroupedExpr(node ast.GroupedExpr) interface{} {
	expr := node.Value.Accept(e)
	return expr.(Object)
}
func (e *Evaluator) VisitPrefixExpr(node ast.PrefixExpr) interface{} {
	expr := node.Right.Accept(e)
	if expr == nil {
		panic("can't evaluate prefix expression")
	}

	switch node.Op {
	case "-":
		if expr, ok := expr.(*NumObject); ok {
			return &NumObject{Value: -expr.Value}
		} else {
			e.Errors = append(e.Errors, errors.New("Operand must be a number."))
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
func (e *Evaluator) VisitInfixExpr(node ast.InfixExpr) interface{} {
	left := node.Left.Accept(e)
	right := node.Right.Accept(e)

	switch node.Op {
	case "+":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &NumObject{Value: l.Value + r.Value}
			}
			e.Errors = append(e.Errors, errors.New("Operands must be two numbers or two strings."))
		}

		if l, ok := left.(*StrObject); ok {
			if r, ok := right.(*StrObject); ok {
				return &StrObject{Value: l.Value + r.Value}
			}
			e.Errors = append(e.Errors, errors.New("Operands must be two numbers or two strings."))
		}
	case "-":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &NumObject{Value: l.Value - r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case "*":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &NumObject{Value: l.Value * r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case "/":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &NumObject{Value: l.Value / r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case "<":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &BooleanObject{Value: l.Value < r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case "<=":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &BooleanObject{Value: l.Value <= r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case ">":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &BooleanObject{Value: l.Value > r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case ">=":
		if l, ok := left.(*NumObject); ok {
			if r, ok := right.(*NumObject); ok {
				return &BooleanObject{Value: l.Value >= r.Value}
			}
		}
		e.Errors = append(e.Errors, errors.New("Operands must be numbers."))
	case "==":
		if left == nil && right == nil {
			return &BooleanObject{Value: true}
		} else if left == nil {
			return &BooleanObject{Value: false}
		}

		return &BooleanObject{Value: reflect.DeepEqual(left, right)}
	case "!=":
		if left == nil && right == nil {
			return &BooleanObject{Value: false}
		} else if left == nil {
			return &BooleanObject{Value: true}
		}

		return &BooleanObject{Value: !reflect.DeepEqual(left, right)}
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

func CheckErrors(errs []error) int {
	for _, err := range errs {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	return 70
}
