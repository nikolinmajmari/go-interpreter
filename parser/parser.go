package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	/// errors
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.NextToken()
	p.NextToken()
	return p
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
	default:
		return nil
	}
}

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
