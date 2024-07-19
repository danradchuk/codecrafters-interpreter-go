package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/ast"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
)

type prefixFunc func() ast.Node
type infixFunc func(left ast.Node) ast.Node

const (
	LOWEST = iota // LOWEST is the universal binding power
	EQUALITY
	COMPARISON
	ADDITIVE
	MULTIPLICATIVE
	PREFIX
	PAREN
)

var tokenTypeToBp = map[lexer.TokenType]int{
	lexer.EQUAL_EQUAL:   EQUALITY,
	lexer.BANG_EQUAL:    EQUALITY,
	lexer.GREATER:       COMPARISON,
	lexer.GREATER_EQUAL: COMPARISON,
	lexer.LESS:          COMPARISON,
	lexer.LESS_EQUAL:    COMPARISON,
	lexer.PLUS:          ADDITIVE,
	lexer.MINUS:         ADDITIVE,
	lexer.STAR:          MULTIPLICATIVE,
	lexer.SLASH:         MULTIPLICATIVE,
	lexer.LEFT_PAREN:    PAREN,
}

type Parser struct {
	Errors    []error
	lexer     *lexer.Lexer
	currToken lexer.Token
	peekToken lexer.Token

	prefixOps map[lexer.TokenType]prefixFunc
	infixOps  map[lexer.TokenType]infixFunc
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:     l,
		prefixOps: make(map[lexer.TokenType]prefixFunc),
		infixOps:  make(map[lexer.TokenType]infixFunc),
	}

	// Prefix
	p.prefixOps[lexer.TRUE] = p.parseBool
	p.prefixOps[lexer.FALSE] = p.parseBool
	p.prefixOps[lexer.NIL] = p.parseNil
	p.prefixOps[lexer.NUMBER] = p.parseNum
	p.prefixOps[lexer.STRING] = p.parseString
	p.prefixOps[lexer.LEFT_PAREN] = p.parseGroupedExpr
	p.prefixOps[lexer.MINUS] = p.parsePrefixExpr
	p.prefixOps[lexer.BANG] = p.parsePrefixExpr

	// Infix
	p.infixOps[lexer.MINUS] = p.parseInfixExpr
	p.infixOps[lexer.PLUS] = p.parseInfixExpr
	p.infixOps[lexer.SLASH] = p.parseInfixExpr
	p.infixOps[lexer.STAR] = p.parseInfixExpr
	p.infixOps[lexer.GREATER] = p.parseInfixExpr
	p.infixOps[lexer.GREATER_EQUAL] = p.parseInfixExpr
	p.infixOps[lexer.LESS] = p.parseInfixExpr
	p.infixOps[lexer.LESS_EQUAL] = p.parseInfixExpr
	p.infixOps[lexer.BANG_EQUAL] = p.parseInfixExpr
	p.infixOps[lexer.EQUAL_EQUAL] = p.parseInfixExpr

	// init currToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseExpr(minBp int) ast.Node {
	prefixFunc := p.prefixOps[p.currToken.Type]
	if prefixFunc == nil {
		p.Errors = append(p.Errors, fmt.Errorf(
			"[line %d] Error at '%s': Expect expression.", p.currToken.Line, p.currToken.Lexeme,
		))
		return nil
	}

	lhs := prefixFunc()
	for minBp < p.peekBp() {
		infixFunc := p.infixOps[p.peekToken.Type]
		if infixFunc == nil {
			return lhs
		}

		p.nextToken() // advance to an operator

		lhs = infixFunc(lhs)
	}

	return lhs
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
func (p *Parser) peekBp() int {
	if bp, ok := tokenTypeToBp[p.peekToken.Type]; ok {
		return bp
	}

	return LOWEST
}
func (p *Parser) currBp() int {
	if bp, ok := tokenTypeToBp[p.currToken.Type]; ok {
		return bp
	}

	return LOWEST
}

func (p *Parser) parseBool() ast.Node {
	b, err := strconv.ParseBool(p.currToken.Lexeme)
	if err != nil {
		p.Errors = append(p.Errors, err)
	}

	return ast.BooleanLiteral{
		Token: p.currToken,
		Value: b,
	}
}
func (p *Parser) parseNil() ast.Node {
	return ast.NilLiteral{}
}
func (p *Parser) parseNum() ast.Node {
	num, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		p.Errors = append(p.Errors, err)
	}

	return ast.NumLiteral{
		Token: p.currToken,
		Value: num,
	}
}
func (p *Parser) parseString() ast.Node {
	return ast.StringLiteral{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}
func (p *Parser) parseGroupedExpr() ast.Node {
	expr := ast.GroupedExpr{
		Token: p.currToken,
	}
	p.nextToken() // consume '('

	exp := p.ParseExpr(0)

	if p.peekToken.Type != lexer.RIGHT_PAREN {
		p.Errors = append(p.Errors, fmt.Errorf("Error: Unmatched parentheses."))
		return nil
	} else {
		p.nextToken() // consume ')'
	}
	expr.Value = exp

	return expr
}
func (p *Parser) parsePrefixExpr() ast.Node {
	expr := ast.PrefixExpr{
		Token: p.currToken,
		Op:    p.currToken.Lexeme,
	}
	p.nextToken() // eat an operator token

	expr.Right = p.ParseExpr(PREFIX)

	return expr
}
func (p *Parser) parseInfixExpr(left ast.Node) ast.Node {
	expr := ast.InfixExpr{
		Token: p.currToken,
		Op:    p.currToken.Lexeme,
		Left:  left,
	}

	rbp := p.currBp()
	p.nextToken() // eat an operator token

	expr.Right = p.ParseExpr(rbp)

	return expr
}

func CheckErrors(errs []error) int {
	for _, err := range errs {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	return 65
}
