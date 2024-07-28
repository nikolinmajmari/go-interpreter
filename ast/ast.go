package ast

import (
	"bytes"
	"interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {

		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.Value)
	out.WriteString(" = ")
	if nil != ls.Value {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) statementNode() {}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if nil != rs.ReturnValue {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) statementNode() {}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if nil != es.Expression {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) statementNode() {}
