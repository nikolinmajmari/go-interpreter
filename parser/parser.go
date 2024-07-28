package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LGT
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	/// errors
	errors []string

	/// parser fns
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) registerPrefix(t token.Type, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}
func (p *Parser) registerInfix(t token.Type, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) PeekError(t token.Type) {
	msg := fmt.Sprintf("Expected token type to be %s, got %s instead!", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	i := 0
	for p.currToken.Type != token.EOF {
		i += 1
		stm := p.ParseStatement()
		if stm != nil {
			program.Statements = append(program.Statements, stm)
		}
		p.NextToken()
	}
	return program
}
func (p *Parser) ParseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

// ParseLetStatement / statement parser
func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stm := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// expectPeek has already advanced parser to new token via p.NextToken()
	stm.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	for !p.currentTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stm
}

func (p *Parser) ParseReturnStatement() ast.Statement {
	stm := &ast.ReturnStatement{Token: p.currToken}
	p.NextToken()
	/// skip expression
	for !p.currentTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stm
}

// ParseExpressionStatement / expressions
func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.ParseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken, Value: p.currToken.Literal,
	}
}

// / helpers
func (p *Parser) currentTokenIs(t token.Type) bool {
	return p.currToken.Type == t
}
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	}
	p.PeekError(t)
	return false
}

// / parsing functions
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
)
