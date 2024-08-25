package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

// / parsing functions
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
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

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LGT,
	token.GT:       LGT,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

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
	p.infixParseFns = make(map[token.Type]infixParseFn)
	/// infix expressions
	p.registerInfix(token.PLUS, p.ParseInfixExpression)
	p.registerInfix(token.MINUS, p.ParseInfixExpression)
	p.registerInfix(token.ASTERISK, p.ParseInfixExpression)
	p.registerInfix(token.SLASH, p.ParseInfixExpression)
	p.registerInfix(token.EQ, p.ParseInfixExpression)
	p.registerInfix(token.NEQ, p.ParseInfixExpression)
	p.registerInfix(token.LT, p.ParseInfixExpression)
	p.registerInfix(token.GT, p.ParseInfixExpression)

	// identifier
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// literal
	p.registerPrefix(token.INT, p.ParseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.ParseBoolean)
	p.registerPrefix(token.FALSE, p.ParseBoolean)

	// prefix expressions
	p.registerPrefix(token.BANG, p.ParsePrefixExpression)
	p.registerPrefix(token.MINUS, p.ParsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.ParseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.ParseFunctionLiteral)

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
	p.NextToken()
	stm.Value = p.ParseExpression(LOWEST)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stm
}

func (p *Parser) ParseReturnStatement() ast.Statement {
	stm := &ast.ReturnStatement{Token: p.currToken}
	p.NextToken()
	/// skip expression
	stm.ReturnValue = p.ParseExpression(LOWEST)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stm
}

// ParseExpressionStatement / expressions
func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	//defer UnTrace(Trace("ParseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.ParseExpression(LOWEST)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) noPrefixParserFnError(t token.Type) {
	p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", t))
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	//defer UnTrace(Trace("ParseExpression"))
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParserFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()
	//
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.NextToken()
		leftExp = infix(leftExp)
		//
	}

	return leftExp
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	//defer UnTrace(Trace("ParsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.NextToken()
	expression.Right = p.ParseExpression(PREFIX)
	return expression
}

func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	//defer UnTrace(Trace("ParseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.NextToken()
	expression.Right = p.ParseExpression(precedence)
	return expression
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	//defer UnTrace(Trace("ParseIntegerLiteral"))
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if nil != err {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	literal.Value = value
	return literal
}

func (p *Parser) ParseBoolean() ast.Expression {
	boolean := &ast.Boolean{Token: p.currToken}
	value, err := strconv.ParseBool(p.currToken.Literal)
	if nil != err {
		msg := fmt.Sprintf("Could not parse %q as boolean", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	boolean.Value = value
	return boolean
}

func (p *Parser) ParseGroupedExpression() ast.Expression {
	p.NextToken()
	exp := p.ParseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
	//defer UnTrace(Trace("ParseIdentifier"))
	return &ast.Identifier{
		Token: p.currToken, Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()
	exp.Condition = p.ParseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	exp.Consequence = p.ParseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.NextToken()
		if !p.expectPeek(token.LBRACKET) {
			return nil
		}
		exp.Alternative = p.ParseBlockStatement()
	}
	return exp
}

func (p *Parser) ParseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken, Statements: []ast.Statement{}}
	p.NextToken()
	for !p.currentTokenIs(token.EOF) && !p.currentTokenIs(token.RBRACKET) {
		stm := p.ParseStatement()
		if stm != nil {
			block.Statements = append(block.Statements, stm)
		}
		p.NextToken()
	}
	return block
}

func (p *Parser) ParseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.currToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.ParseParameters()

	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	lit.Body = p.ParseBlockStatement()
	fmt.Printf("%v", lit)
	return lit
}

func (p *Parser) ParseParameters() []*ast.Identifier {
	var params []*ast.Identifier
	if p.peekTokenIs(token.RPAREN) {
		p.NextToken()
		return params
	}
	p.NextToken()
	param := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		param := &ast.Identifier{
			Token: p.currToken,
			Value: p.currToken.Literal,
		}
		params = append(params, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return params
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

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}
