package ast

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
)

type Visitor interface {
	VisitBoolean(node BooleanLiteral)
	VisitNil(node NilLiteral)
	VisitNum(node NumLiteral)
	VisitString(node StringLiteral)
	VisitGroupedExpr(node GroupedExpr)
	VisitPrefixExpr(node PrefixExpr)
	VisitInfixExpr(node InfixExpr)
}

type Node interface {
	Type() string
	String() string
	Accept(visitor Visitor)
}

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (n BooleanLiteral) Type() string {
	return "BOOLEAN"
}
func (n BooleanLiteral) String() string {
	return fmt.Sprintf("%t", n.Value)
}
func (n BooleanLiteral) Accept(visitor Visitor) { visitor.VisitBoolean(n) }

type NilLiteral struct{}

func (n NilLiteral) Type() string {
	return "NIL"
}
func (n NilLiteral) String() string         { return fmt.Sprintf("%s", "nil") }
func (n NilLiteral) Accept(visitor Visitor) { visitor.VisitNil(n) }

type NumLiteral struct {
	Token lexer.Token
	Value float64
}

func (n NumLiteral) Type() string {
	return "NUMBER"
}
func (n NumLiteral) String() string {
	return trailZeroes(fmt.Sprintf("%f", n.Value))
}
func (n NumLiteral) Accept(visitor Visitor) { visitor.VisitNum(n) }

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (n StringLiteral) Type() string           { return "STRING" }
func (n StringLiteral) String() string         { return fmt.Sprintf("%s", n.Value) }
func (n StringLiteral) Accept(visitor Visitor) { visitor.VisitString(n) }

type GroupedExpr struct {
	Token lexer.Token
	Value Node
}

func (n GroupedExpr) Type() string           { return "GROUPED_EXPR" }
func (n GroupedExpr) String() string         { return parenthesize("group", n.Value) }
func (n GroupedExpr) Accept(visitor Visitor) { visitor.VisitGroupedExpr(n) }

type PrefixExpr struct {
	Token lexer.Token
	Op    string
	Right Node
}

func (n PrefixExpr) Type() string           { return "PREFIX_EXPR" }
func (n PrefixExpr) String() string         { return parenthesize(n.Op, n.Right) }
func (n PrefixExpr) Accept(visitor Visitor) { visitor.VisitPrefixExpr(n) }

type InfixExpr struct {
	Token lexer.Token
	Left  Node
	Op    string
	Right Node
}

func (n InfixExpr) Type() string           { return "INFIX_EXPR" }
func (n InfixExpr) String() string         { return parenthesize(n.Op, n.Left, n.Right) }
func (n InfixExpr) Accept(visitor Visitor) { visitor.VisitInfixExpr(n) }

func parenthesize(op string, expr ...Node) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(op)

	for _, e := range expr {
		sb.WriteString(" ")
		sb.WriteString(e.String())
	}

	sb.WriteString(")")

	return sb.String()
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
