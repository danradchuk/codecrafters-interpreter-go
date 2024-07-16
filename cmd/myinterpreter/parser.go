package main

import (
	"fmt"
	"os"
	"strconv"
)

type prefixFunc func() ASTNode
type infixFunc func(left ASTNode) ASTNode

const (
	LOWEST = iota // LOWEST is the universal binding power
	EQUALITY
	COMPARISON
	ADDITIVE
	MULTIPLICATIVE
	PREFIX
	PAREN
)

var tokenTypeToBp = map[TokenType]int{
	EQUAL:         EQUALITY,
	BANG_EQUAL:    EQUALITY,
	GREATER:       COMPARISON,
	GREATER_EQUAL: COMPARISON,
	LESS:          COMPARISON,
	LESS_EQUAL:    COMPARISON,
	PLUS:          ADDITIVE,
	MINUS:         ADDITIVE,
	STAR:          MULTIPLICATIVE,
	SLASH:         MULTIPLICATIVE,
	LEFT_PAREN:    PAREN,
}

type Parser struct {
	lexer     *Lexer
	errs      []error
	currToken Token
	peekToken Token

	prefixOps map[TokenType]prefixFunc
	infixOps  map[TokenType]infixFunc
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:     lexer,
		prefixOps: make(map[TokenType]prefixFunc),
		infixOps:  make(map[TokenType]infixFunc),
	}

	// Prefix
	p.prefixOps[TRUE] = p.parseBool
	p.prefixOps[FALSE] = p.parseBool
	p.prefixOps[NIL] = p.parseNil
	p.prefixOps[NUMBER] = p.parseNum
	p.prefixOps[STRING] = p.parseString
	p.prefixOps[LEFT_PAREN] = p.parseGroupedExpr
	p.prefixOps[MINUS] = p.parsePrefixExpr
	p.prefixOps[BANG] = p.parsePrefixExpr

	// Infix
	p.infixOps[MINUS] = p.parseInfixExpr
	p.infixOps[PLUS] = p.parseInfixExpr
	p.infixOps[SLASH] = p.parseInfixExpr
	p.infixOps[STAR] = p.parseInfixExpr
	p.infixOps[GREATER] = p.parseInfixExpr
	p.infixOps[GREATER_EQUAL] = p.parseInfixExpr
	p.infixOps[LESS] = p.parseInfixExpr
	p.infixOps[LESS_EQUAL] = p.parseInfixExpr
	p.infixOps[BANG_EQUAL] = p.parseInfixExpr
	p.infixOps[EQUAL_EQUAL] = p.parseInfixExpr

	// init currToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseExpr(minBp int) ASTNode {
	prefixFunc := p.prefixOps[p.currToken.Type]
	if prefixFunc == nil {
		// TODO add error
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

func (p *Parser) parseBool() ASTNode {
	b, err := strconv.ParseBool(p.currToken.Lexeme)
	if err != nil {
		p.errs = append(p.errs, err)
	}

	return BooleanLiteral{
		Token: p.currToken,
		Value: b,
	}
}
func (p *Parser) parseNil() ASTNode {
	return NilLiteral{}
}
func (p *Parser) parseNum() ASTNode {
	num, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		p.errs = append(p.errs, err)
	}

	return NumLiteral{
		Token: p.currToken,
		Value: num,
	}
}
func (p *Parser) parseString() ASTNode {
	return StringLiteral{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}
func (p *Parser) parseGroupedExpr() ASTNode {
	expr := GroupedExpr{
		Token: p.currToken,
	}
	p.nextToken()

	exp := p.ParseExpr(0)

	if p.peekToken.Type != RIGHT_PAREN {
		p.errs = append(p.errs, fmt.Errorf("Error: Unmatched parentheses."))
		return nil
	} else {
		p.nextToken()
	}
	expr.Value = exp

	return expr
}
func (p *Parser) parsePrefixExpr() ASTNode {
	expr := PrefixExpr{
		Token: p.currToken,
		Op:    p.currToken.Lexeme,
	}
	p.nextToken() // eat an operator token

	expr.Right = p.ParseExpr(PREFIX)

	return expr
}
func (p *Parser) parseInfixExpr(left ASTNode) ASTNode {
	expr := InfixExpr{
		Token: p.currToken,
		Op:    p.currToken.Lexeme,
		Left:  left,
	}

	rbp := p.currBp()
	p.nextToken() // eat an operator token

	expr.Right = p.ParseExpr(rbp)

	return expr
}

func CheckParserErrors(errs []error) int {
	for _, err := range errs {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	return 65
}
