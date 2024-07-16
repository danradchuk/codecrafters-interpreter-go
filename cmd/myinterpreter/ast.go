package main

import (
	"fmt"
	"strings"
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

type ASTNode interface {
	Type() string
	String() string
	Accept(visitor Visitor)
}

type BooleanLiteral struct {
	Token Token
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
	Token Token
	Value float64
}

func (n NumLiteral) Type() string {
	return "NUMBER"
}
func (n NumLiteral) String() string         { return fmt.Sprintf("%f", n.Value) }
func (n NumLiteral) Accept(visitor Visitor) { visitor.VisitNum(n) }

type StringLiteral struct {
	Token Token
	Value string
}

func (n StringLiteral) Type() string           { return "STRING" }
func (n StringLiteral) String() string         { return fmt.Sprintf("%s", n.Value) }
func (n StringLiteral) Accept(visitor Visitor) { visitor.VisitString(n) }

type GroupedExpr struct {
	Token Token
	Value ASTNode
}

func (n GroupedExpr) Type() string           { return "GROUPED_EXPR" }
func (n GroupedExpr) String() string         { return parenthesize("group", n.Value) }
func (n GroupedExpr) Accept(visitor Visitor) { visitor.VisitGroupedExpr(n) }

type PrefixExpr struct {
	Token Token
	Op    string
	Right ASTNode
}

func (n PrefixExpr) Type() string           { return "PREFIX_EXPR" }
func (n PrefixExpr) String() string         { return parenthesize(n.Op, n.Right) }
func (n PrefixExpr) Accept(visitor Visitor) { visitor.VisitPrefixExpr(n) }

type InfixExpr struct {
	Token Token
	Left  ASTNode
	Op    string
	Right ASTNode
}

func (n InfixExpr) Type() string           { return "INFIX_EXPR" }
func (n InfixExpr) String() string         { return parenthesize(n.Op, n.Left, n.Right) }
func (n InfixExpr) Accept(visitor Visitor) { visitor.VisitInfixExpr(n) }

func parenthesize(op string, expr ...ASTNode) string {
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
